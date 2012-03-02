/*
Rating system designed to be used in VoIP Carriers World
Copyright (C) 2012  Radu Ioan Fericean

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/
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
	storageGetter                      StorageGetter
	userBudget                         *UserBudget
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
Gets and caches the user budget information.
*/
func (cd *CallDescriptor) getUserBudget() (ub *UserBudget, err error) {
	if cd.userBudget == nil {
		cd.userBudget, err = cd.storageGetter.GetUserBudget(cd.Subject)
	}
	return cd.userBudget, err
}

/*
Exported method to set the storage getter.
*/
func (cd *CallDescriptor) SetStorageGetter(sg StorageGetter) {
	cd.storageGetter = sg
}

/*
Restores the activation periods for the specified prefix from storage.
*/
func (cd *CallDescriptor) SearchStorageForPrefix() (destPrefix string, err error) {
	cd.ActivationPeriods = make([]*ActivationPeriod, 0)
	base := fmt.Sprintf("%s:%s:", cd.CstmId, cd.Subject)
	destPrefix = cd.DestinationPrefix
	key := base + destPrefix
	values, err := cd.storageGetter.GetActivationPeriods(key)
	//get for a smaller prefix if the orignal one was not found	
	for i := len(cd.DestinationPrefix); err != nil &&
		i >= MinPrefixLength; values, err = cd.storageGetter.GetActivationPeriods(key) {
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
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		userBudget.mux.RLock()
		_, bucketList := userBudget.getSecondsForPrefix(cd.storageGetter, cd.DestinationPrefix)
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
	destPrefix, err := cd.SearchStorageForPrefix()

	timespans := cd.splitInTimeSpans()

	cost := 0.0
	connectionFee := 0.0
	for i, ts := range timespans {
		if i == 0 && ts.MinuteInfo == nil && ts.Interval != nil {
			connectionFee = ts.Interval.ConnectFee
		}
		cost += ts.GetCost(cd)
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
	// TODO: remove this method if if not still used
	_, err = cd.SearchStorageForPrefix()
	now := time.Now()
	oneSecond, _ := time.ParseDuration("1s")
	ts := &TimeSpan{TimeStart: now, TimeEnd: now.Add(oneSecond)}
	timespans := cd.splitTimeSpan(ts)

	if len(timespans) > 0 {
		cost = round(timespans[0].GetCost(cd), 3)
	}
	return
}

/*
Returns the aproximate max allowed session for user budget. It will try the max amount received in the call descriptor 
and will decrease it by 10% for nine times. So if the user has little credit it will still allow 10% of the initial amount.
If the user has no credit then it will return 0.
*/
func (cd *CallDescriptor) GetMaxSessionTime() (seconds float64, err error) {
	_, err = cd.SearchStorageForPrefix()
	now := time.Now()
	availableCredit, availableSeconds := 0.0, 0.0
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		userBudget.mux.RLock()
		availableCredit = userBudget.Credit
		availableSeconds, _ = userBudget.getSecondsForPrefix(cd.storageGetter, cd.DestinationPrefix)
		userBudget.mux.RUnlock()
	} else {
		return cd.Amount, err
	}
	// check for zero budget	
	if availableCredit == 0 {
		return availableSeconds, nil
	}

	maxSessionSeconds := cd.Amount
	for i := 0; i < 10; i++ {
		maxDuration, _ := time.ParseDuration(fmt.Sprintf("%vs", maxSessionSeconds))
		ts := &TimeSpan{TimeStart: now, TimeEnd: now.Add(maxDuration)}
		timespans := cd.splitTimeSpan(ts)

		cost := 0.0
		for i, ts := range timespans {
			if i == 0 && ts.MinuteInfo == nil && ts.Interval != nil {
				cost += ts.Interval.ConnectFee
			}
			cost += ts.GetCost(cd)
		}
		if cost < availableCredit {
			return maxSessionSeconds, nil
		} else { //decrease the period by 10% and try again
			maxSessionSeconds -= cd.Amount * 0.1
		}
	}
	return 0, nil
}

/*
Interface method used to add/substract an amount of cents from user's money budget.
The amount filed has to be filled in call descriptor.
*/
func (cd *CallDescriptor) DebitCents() (left float64, err error) {
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		return userBudget.debitMoneyBudget(cd.storageGetter, cd.Amount), nil
	}
	return 0.0, err
}

/*
Interface method used to add/substract an amount of units from user's sms budget.
The amount filed has to be filled in call descriptor.
*/
func (cd *CallDescriptor) DebitSMS() (left float64, err error) {
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		return userBudget.debitSMSBuget(cd.storageGetter, cd.Amount)
	}
	return 0, err
}

/*
Interface method used to add/substract an amount of seconds from user's minutes budget.
The amount filed has to be filled in call descriptor.
*/
func (cd *CallDescriptor) DebitSeconds() (err error) {
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		return userBudget.debitMinutesBudget(cd.storageGetter, cd.Amount, cd.DestinationPrefix)
	}
	return err
}

/*
Interface method used to add an amount to the accumulated placed call seconds
to be used for volume discount.
The amount filed has to be filled in call descriptor.
*/
func (cd *CallDescriptor) AddVolumeDiscountSeconds() (err error) {
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		return userBudget.addVolumeDiscountSeconds(cd.storageGetter, cd.Amount)
	}
	return err
}

/*
Resets the accumulated volume discount seconds (to zero).
*/
func (cd *CallDescriptor) ResetVolumeDiscountSeconds() (err error) {
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		return userBudget.resetVolumeDiscountSeconds(cd.storageGetter)
	}
	return err
}

/*
Adds the specified amount of seconds to the recived call seconds. When the threshold specified
in the user's tariff plan is reached then the recived call budget is reseted and the bonus
specified in the tariff plan is applyed.
The amount filed has to be filled in call descriptor.
*/
func (cd *CallDescriptor) AddRecievedCallSeconds() (err error) {
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		return userBudget.addReceivedCallSeconds(cd.storageGetter, cd.Amount)
	}
	return err
}

/*
Resets user budgets value to the amounts specified in the tariff plan.
*/
func (cd *CallDescriptor) ResetUserBudget() (err error) {
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		return userBudget.resetUserBudget(cd.storageGetter)
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
