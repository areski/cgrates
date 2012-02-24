package timespans

import (
	"fmt"
	// "log"
	"math"
	"time"
)

const (
	// the minimum length for a destination prefix to be matched.
	MinPrefixLength = 2
)

/*
Utility function for rounding a float to a certain number of decimals (not present in math).
*/
func round(val float64, prec int) float64 {

	var rounder float64
	intermed := val * math.Pow(10, float64(prec))

	if val >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / math.Pow(10, float64(prec))
}

/*
The input stucture that contains call information.
*/
type CallDescriptor struct {
	TOR                                int
	CstmId, Subject, DestinationPrefix string
	TimeStart, TimeEnd                 time.Time
	Amount                             float64
	ActivationPeriods                  []*ActivationPeriod
	StorageGetter                      StorageGetter
}

/*
Adds an activation period that applyes to current call descriptor.
*/
func (cd *CallDescriptor) AddActivationPeriod(aps ...*ActivationPeriod) {
	for _, ap := range aps {
		cd.ActivationPeriods = append(cd.ActivationPeriods, ap)
	}
}

/*
Restores the activation periods from storage.
*/
func (cd *CallDescriptor) RestoreFromStorage() (destPrefix string, err error) {
	cd.ActivationPeriods = make([]*ActivationPeriod, 0)
	base := fmt.Sprintf("%s:%s:", cd.CstmId, cd.Subject)
	destPrefix = cd.DestinationPrefix
	key := base + destPrefix
	values, err := cd.StorageGetter.GetActivationPeriods(key)
	//get for a smaller prefix if the orignal one was not found

	for i := len(cd.DestinationPrefix); err != nil && i >= MinPrefixLength; values, err = cd.StorageGetter.GetActivationPeriods(key) {
		i--
		destPrefix = cd.DestinationPrefix[:i]
		key = base + destPrefix
	}
	//load the activation preriods
	if err == nil {
		cd.ActivationPeriods = values
	}
	return
}

/*
Constructs the key for the storage lookup.
The prefixLen is limiting the length of the destination prefix.
*/
func (cd *CallDescriptor) GetKey() string {
	return fmt.Sprintf("%s:%s:%s", cd.CstmId, cd.Subject, cd.DestinationPrefix)
}

/*
Splits the call descriptor timespan into sub time spans according to the activation periods intervals.
*/
func (cd *CallDescriptor) splitInTimeSpans() (timespans []*TimeSpan) {
	return cd.splitTimeSpan(&TimeSpan{TimeStart: cd.TimeStart, TimeEnd: cd.TimeEnd})
}

/*
Splits the received timespan into sub time spans according to the activation periods intervals.
*/
func (cd *CallDescriptor) splitTimeSpan(firstSpan *TimeSpan) (timespans []*TimeSpan) {
	timespans = append(timespans, firstSpan)
	// split on (free) minute buckets	
	if userBudget, err := cd.StorageGetter.GetUserBudget(cd.Subject); err == nil && userBudget != nil {
		userBudget.mux.RLock()
		_, bucketList := userBudget.getSecondsForPrefix(cd.StorageGetter, cd.DestinationPrefix)
		for _, mb := range bucketList {
			for i := 0; i < len(timespans); i++ {
				if timespans[i].MinuteInfo != nil {
					continue
				}
				newTs := timespans[i].SplitByMinuteBucket(mb)
				if newTs != nil {
					timespans = append(timespans, newTs)
					firstSpan = newTs // we move the firstspan to the newly created one for further spliting
					break
				}
			}
		}
		userBudget.mux.RUnlock()
	}

	if firstSpan.MinuteInfo != nil {
		return // all the timespans are on minutes
	}
	if len(cd.ActivationPeriods) == 0 {
		return
	}

	firstSpan.ActivationPeriod = cd.ActivationPeriods[0]

	// split on activation periods
	afterStart, afterEnd := false, false //optimization for multiple activation periods
	for _, ap := range cd.ActivationPeriods {
		if !afterStart && !afterEnd && ap.ActivationTime.Before(cd.TimeStart) {
			firstSpan.ActivationPeriod = ap
		} else {
			afterStart = true
			for i := 0; i < len(timespans); i++ {
				if timespans[i].MinuteInfo != nil {
					continue
				}
				newTs := timespans[i].SplitByActivationPeriod(ap)
				if newTs != nil {
					timespans = append(timespans, newTs)
				} else {
					afterEnd = true
					break
				}
			}
		}
	}
	// split on price intervals
	for i := 0; i < len(timespans); i++ {
		if timespans[i].MinuteInfo != nil {
			continue
		}
		ap := timespans[i].ActivationPeriod
		//timespans[i].ActivationPeriod = nil
		for _, interval := range ap.Intervals {
			newTs := timespans[i].SplitByInterval(interval)
			if newTs != nil {
				newTs.ActivationPeriod = ap
				timespans = append(timespans, newTs)
			}
		}
	}
	return
}

/*
Creates a CallCost structure with the cost nformation calculated for the received CallDescriptor.
*/
func (cd *CallDescriptor) GetCost() (*CallCost, error) {
	destPrefix, err := cd.RestoreFromStorage()

	timespans := cd.splitInTimeSpans()

	cost := 0.0
	connectionFee := 0.0
	for i, ts := range timespans {
		if ts.MinuteInfo == nil && i == 0 {
			connectionFee = ts.Interval.ConnectFee
		}
		cost += ts.GetCost()
	}
	cc := &CallCost{TOR: cd.TOR,
		CstmId:            cd.CstmId,
		Subject:           cd.Subject,
		DestinationPrefix: destPrefix,
		Cost:              cost,
		ConnectFee:        connectionFee,
		Timespans:         timespans}

	return cc, err
}

/*
Returns the cost of a second in the present time conditions.
*/
func (cd *CallDescriptor) getPresentSecondCost() (cost float64, err error) {
	_, err = cd.RestoreFromStorage()
	now := time.Now()
	oneSecond, _ := time.ParseDuration("1s")
	ts := &TimeSpan{TimeStart: now, TimeEnd: now.Add(oneSecond)}
	timespans := cd.splitTimeSpan(ts)

	if len(timespans) > 0 {
		cost = round(timespans[0].GetCost(), 3)
	}
	return
}

/*
Returns the cost of a second in the present time conditions.
*/
func (cd *CallDescriptor) GetMaxSessionTime(maxSessionSeconds int) (seconds int, err error) {
	_, err = cd.RestoreFromStorage()
	now := time.Now()
	availableCredit := 0.0
	if userBudget, err := cd.StorageGetter.GetUserBudget(cd.Subject); err == nil && userBudget != nil {
		availableCredit = userBudget.Credit
	} else {
		return maxSessionSeconds, err
	}
	orig_maxSessionSeconds := maxSessionSeconds
	for i := 0; i < 10; i++ {
		maxDuration, _ := time.ParseDuration(fmt.Sprintf("%vs", maxSessionSeconds))
		ts := &TimeSpan{TimeStart: now, TimeEnd: now.Add(maxDuration)}
		timespans := cd.splitTimeSpan(ts)

		cost := 0.0
		for i, ts := range timespans {
			if ts.MinuteInfo == nil && i == 0 {
				cost += ts.Interval.ConnectFee
			}
			cost += ts.GetCost()
		}
		if cost < availableCredit {
			return maxSessionSeconds, nil
		} else { //decrease the period by 10% and try again
			maxSessionSeconds -= int(float64(orig_maxSessionSeconds) * 0.1)
		}
	}
	return 0, nil
}

func (cd *CallDescriptor) DebitCents() (left float64, err error) {
	if userBudget, err := cd.StorageGetter.GetUserBudget(cd.Subject); err == nil && userBudget != nil {
		return userBudget.debitMoneyBudget(cd.StorageGetter, cd.Amount), nil
	}
	return 0.0, err
}

func (cd *CallDescriptor) DebitSMS() (left float64, err error) {
	if userBudget, err := cd.StorageGetter.GetUserBudget(cd.Subject); err == nil && userBudget != nil {
		return userBudget.debitSMSBuget(cd.StorageGetter, cd.Amount)
	}
	return 0, err
}

func (cd *CallDescriptor) DebitSeconds() (err error) {
	if userBudget, err := cd.StorageGetter.GetUserBudget(cd.Subject); err == nil && userBudget != nil {
		return userBudget.debitMinutesBudget(cd.StorageGetter, cd.Amount, cd.DestinationPrefix)
	}
	return err
}

/*
The output structure that will be returned with the call cost information.
*/
type CallCost struct {
	TOR                                int
	CstmId, Subject, DestinationPrefix string
	Cost, ConnectFee                   float64
	Timespans                          []*TimeSpan
}