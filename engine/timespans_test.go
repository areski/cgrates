/*
Rating system designed to be used in VoIP Carriers World
Copyright (C) 2013 ITsysCOM

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

package engine

import (
	//"github.com/cgrates/cgrates/utils"
	"testing"
	"time"
)

func TestRightMargin(t *testing.T) {
	i := &RateInterval{WeekDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}}
	t1 := time.Date(2012, time.February, 3, 23, 45, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 4, 0, 10, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	if ts.TimeStart != t1 || ts.TimeEnd != time.Date(2012, time.February, 3, 23, 59, 59, 0, time.UTC) {
		t.Error("Incorrect first half", ts)
	}
	if nts.TimeStart != time.Date(2012, time.February, 3, 23, 59, 59, 0, time.UTC) || nts.TimeEnd != t2 {
		t.Error("Incorrect second half", nts)
	}
	if ts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}

	if ts.GetDuration().Seconds() != 15*60-1 || nts.GetDuration().Seconds() != 10*60+1 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), ts.GetDuration().Seconds())
	}

	if ts.GetDuration().Seconds()+nts.GetDuration().Seconds() != oldDuration.Seconds() {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
}

func TestRightHourMargin(t *testing.T) {
	i := &RateInterval{WeekDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}, EndTime: "17:59:00"}
	t1 := time.Date(2012, time.February, 3, 17, 30, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 3, 18, 00, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	if ts.TimeStart != t1 || ts.TimeEnd != time.Date(2012, time.February, 3, 17, 59, 00, 0, time.UTC) {
		t.Error("Incorrect first half", ts)
	}
	if nts.TimeStart != time.Date(2012, time.February, 3, 17, 59, 00, 0, time.UTC) || nts.TimeEnd != t2 {
		t.Error("Incorrect second half", nts)
	}
	if ts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}

	if ts.GetDuration().Seconds() != 29*60 || nts.GetDuration().Seconds() != 1*60 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), nts.GetDuration().Seconds())
	}
	if ts.GetDuration().Seconds()+nts.GetDuration().Seconds() != oldDuration.Seconds() {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
}

func TestLeftMargin(t *testing.T) {
	i := &RateInterval{WeekDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}}
	t1 := time.Date(2012, time.February, 5, 23, 45, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 6, 0, 10, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	if ts.TimeStart != t1 || ts.TimeEnd != time.Date(2012, time.February, 6, 0, 0, 0, 0, time.UTC) {
		t.Error("Incorrect first half", ts)
	}
	if nts.TimeStart != time.Date(2012, time.February, 6, 0, 0, 0, 0, time.UTC) || nts.TimeEnd != t2 {
		t.Error("Incorrect second half", nts)
	}
	if nts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}
	if ts.GetDuration().Seconds() != 15*60 || nts.GetDuration().Seconds() != 10*60 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), nts.GetDuration().Seconds())
	}
	if ts.GetDuration().Seconds()+nts.GetDuration().Seconds() != oldDuration.Seconds() {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
}

func TestLeftHourMargin(t *testing.T) {
	i := &RateInterval{Months: Months{time.December}, MonthDays: MonthDays{1}, StartTime: "09:00:00"}
	t1 := time.Date(2012, time.December, 1, 8, 45, 0, 0, time.UTC)
	t2 := time.Date(2012, time.December, 1, 9, 20, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	if ts.TimeStart != t1 || ts.TimeEnd != time.Date(2012, time.December, 1, 9, 0, 0, 0, time.UTC) {
		t.Error("Incorrect first half", ts)
	}
	if nts.TimeStart != time.Date(2012, time.December, 1, 9, 0, 0, 0, time.UTC) || nts.TimeEnd != t2 {
		t.Error("Incorrect second half", nts)
	}
	if nts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}
	if ts.GetDuration().Seconds() != 15*60 || nts.GetDuration().Seconds() != 20*60 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), nts.GetDuration().Seconds())
	}
	if ts.GetDuration().Seconds()+nts.GetDuration().Seconds() != oldDuration.Seconds() {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
}

func TestEnclosingMargin(t *testing.T) {
	i := &RateInterval{WeekDays: []time.Weekday{time.Sunday}}
	t1 := time.Date(2012, time.February, 5, 17, 45, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 5, 18, 10, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2}
	nts := ts.SplitByRateInterval(i)
	if ts.TimeStart != t1 || ts.TimeEnd != t2 || nts != nil {
		t.Error("Incorrect enclosing", ts)
	}
	if ts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}
}

func TestOutsideMargin(t *testing.T) {
	i := &RateInterval{WeekDays: []time.Weekday{time.Monday}}
	t1 := time.Date(2012, time.February, 5, 17, 45, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 5, 18, 10, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2}
	result := ts.SplitByRateInterval(i)
	if result != nil {
		t.Error("RateInterval not split correctly")
	}
}

func TestContains(t *testing.T) {
	t1 := time.Date(2012, time.February, 5, 17, 45, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 5, 17, 55, 0, 0, time.UTC)
	t3 := time.Date(2012, time.February, 5, 17, 50, 0, 0, time.UTC)
	ts := TimeSpan{TimeStart: t1, TimeEnd: t2}
	if ts.Contains(t1) {
		t.Error("It should NOT contain ", t1)
	}
	if ts.Contains(t2) {
		t.Error("It should NOT contain ", t1)
	}
	if !ts.Contains(t3) {
		t.Error("It should contain ", t3)
	}
}

func TestSplitByActivationTime(t *testing.T) {
	t1 := time.Date(2012, time.February, 5, 17, 45, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 5, 17, 55, 0, 0, time.UTC)
	t3 := time.Date(2012, time.February, 5, 17, 50, 0, 0, time.UTC)
	ts := TimeSpan{TimeStart: t1, TimeEnd: t2}
	ap1 := &RatingPlan{ActivationTime: t1}
	ap2 := &RatingPlan{ActivationTime: t2}
	ap3 := &RatingPlan{ActivationTime: t3}

	if ts.SplitByRatingPlan(ap1) != nil {
		t.Error("Error spliting on left margin")
	}
	if ts.SplitByRatingPlan(ap2) != nil {
		t.Error("Error spliting on right margin")
	}
	result := ts.SplitByRatingPlan(ap3)
	if result.TimeStart != t3 || result.TimeEnd != t2 {
		t.Error("Error spliting on interior")
	}
}

func TestTimespanGetCost(t *testing.T) {
	t1 := time.Date(2012, time.February, 5, 17, 45, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 5, 17, 55, 0, 0, time.UTC)
	ts1 := TimeSpan{TimeStart: t1, TimeEnd: t2}
	if ts1.getCost() != 0 {
		t.Error("No interval and still kicking")
	}
	ts1.SetRateInterval(&RateInterval{Rates: RateGroups{&Rate{0, 1.0, 1 * time.Second, 1 * time.Second}}})
	if ts1.getCost() != 600 {
		t.Error("Expected 10 got ", ts1.Cost)
	}
	ts1.RateInterval = nil
	ts1.SetRateInterval(&RateInterval{Rates: RateGroups{&Rate{0, 1.0, 1 * time.Second, 60 * time.Second}}})
	if ts1.getCost() != 10 {
		t.Error("Expected 6000 got ", ts1.Cost)
	}
}

func TestSetRateInterval(t *testing.T) {
	i1 := &RateInterval{Rates: RateGroups{&Rate{0, 1.0, 1 * time.Second, 1 * time.Second}}}
	ts1 := TimeSpan{RateInterval: i1}
	i2 := &RateInterval{Rates: RateGroups{&Rate{0, 2.0, 1 * time.Second, 1 * time.Second}}}
	ts1.SetRateInterval(i2)
	if ts1.RateInterval != i1 {
		t.Error("Smaller price interval should win")
	}
	i2.Weight = 1
	ts1.SetRateInterval(i2)
	if ts1.RateInterval != i2 {
		t.Error("Bigger ponder interval should win")
	}
}

func TestTimespanSplitGroupedRates(t *testing.T) {
	i := &RateInterval{
		EndTime: "17:59:00",
		Rates:   RateGroups{&Rate{0, 2, 1 * time.Second, 1 * time.Second}, &Rate{900 * time.Second, 1, 1 * time.Second, 1 * time.Second}},
	}
	t1 := time.Date(2012, time.February, 3, 17, 30, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 3, 18, 00, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2, CallDuration: 1800 * time.Second}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	splitTime := time.Date(2012, time.February, 3, 17, 45, 00, 0, time.UTC)
	if ts.TimeStart != t1 || ts.TimeEnd != splitTime {
		t.Error("Incorrect first half", ts.TimeStart, ts.TimeEnd)
	}
	if nts.TimeStart != splitTime || nts.TimeEnd != t2 {
		t.Error("Incorrect second half", nts)
	}
	if ts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}
	c1 := ts.RateInterval.GetCost(ts.GetDuration(), ts.GetGroupStart())
	c2 := nts.RateInterval.GetCost(nts.GetDuration(), nts.GetGroupStart())
	if c1 != 1800 || c2 != 900 {
		t.Error("Wrong costs: ", c1, c2)
	}

	if ts.GetDuration().Seconds() != 15*60 || nts.GetDuration().Seconds() != 15*60 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), nts.GetDuration().Seconds())
	}
	if ts.GetDuration().Seconds()+nts.GetDuration().Seconds() != oldDuration.Seconds() {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
}

func TestTimespanSplitGroupedRatesIncrements(t *testing.T) {
	i := &RateInterval{
		EndTime: "17:59:00",
		Rates: RateGroups{
			&Rate{
				GroupIntervalStart: 0,
				Value:              2,
				RateIncrement:      time.Second,
				RateUnit:           time.Second},
			&Rate{
				GroupIntervalStart: 30 * time.Second,
				Value:              1,
				RateIncrement:      time.Minute,
				RateUnit:           time.Second,
			}},
	}
	t1 := time.Date(2012, time.February, 3, 17, 30, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 3, 17, 31, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2, CallDuration: 60 * time.Second}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	cd := &CallDescriptor{}
	timespans := cd.roundTimeSpansToIncrement([]*TimeSpan{ts, nts})
	if len(timespans) != 2 {
		t.Error("Error rounding timespans: ", timespans)
	}
	ts = timespans[0]
	nts = timespans[1]
	splitTime := time.Date(2012, time.February, 3, 17, 30, 30, 0, time.UTC)
	if ts.TimeStart != t1 || ts.TimeEnd != splitTime {
		t.Error("Incorrect first half", ts)
	}
	t3 := time.Date(2012, time.February, 3, 17, 31, 30, 0, time.UTC)
	if nts.TimeStart != splitTime || nts.TimeEnd != t3 {
		t.Error("Incorrect second half", nts.TimeStart, nts.TimeEnd)
	}
	if ts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}
	c1 := ts.RateInterval.GetCost(ts.GetDuration(), ts.GetGroupStart())
	c2 := nts.RateInterval.GetCost(nts.GetDuration(), nts.GetGroupStart())
	if c1 != 60 || c2 != 60 {
		t.Error("Wrong costs: ", c1, c2)
	}

	if ts.GetDuration().Seconds() != 0.5*60 || nts.GetDuration().Seconds() != 1*60 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), nts.GetDuration().Seconds())
	}
	if ts.GetDuration()+nts.GetDuration() != oldDuration+30*time.Second {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
}

func TestTimespanSplitRightHourMarginBeforeGroup(t *testing.T) {
	i := &RateInterval{
		EndTime: "17:00:30",
		Rates:   RateGroups{&Rate{0, 2, 1 * time.Second, 1 * time.Second}, &Rate{60 * time.Second, 1, 60 * time.Second, 1 * time.Second}},
	}
	t1 := time.Date(2012, time.February, 3, 17, 00, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 3, 17, 01, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	splitTime := time.Date(2012, time.February, 3, 17, 00, 30, 0, time.UTC)
	if ts.TimeStart != t1 || ts.TimeEnd != splitTime {
		t.Error("Incorrect first half", ts)
	}
	if nts.TimeStart != splitTime || nts.TimeEnd != t2 {
		t.Error("Incorrect second half", nts)
	}
	if ts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}

	if ts.GetDuration().Seconds() != 30 || nts.GetDuration().Seconds() != 30 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), nts.GetDuration().Seconds())
	}
	if ts.GetDuration().Seconds()+nts.GetDuration().Seconds() != oldDuration.Seconds() {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
	nnts := nts.SplitByRateInterval(i)
	if nnts != nil {
		t.Error("Bad new split", nnts)
	}
}

func TestTimespanSplitGroupSecondSplit(t *testing.T) {
	i := &RateInterval{
		EndTime: "17:03:30",
		Rates:   RateGroups{&Rate{0, 2, 1 * time.Second, 1 * time.Second}, &Rate{60 * time.Second, 1, 1 * time.Second, 1 * time.Second}},
	}
	t1 := time.Date(2012, time.February, 3, 17, 00, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 3, 17, 04, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2, CallDuration: 240 * time.Second}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	splitTime := time.Date(2012, time.February, 3, 17, 01, 00, 0, time.UTC)
	if ts.TimeStart != t1 || ts.TimeEnd != splitTime {
		t.Error("Incorrect first half", nts)
	}
	if nts.TimeStart != splitTime || nts.TimeEnd != t2 {
		t.Error("Incorrect second half", nts)
	}
	if ts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}

	if ts.GetDuration().Seconds() != 60 || nts.GetDuration().Seconds() != 180 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), nts.GetDuration().Seconds())
	}
	if ts.GetDuration().Seconds()+nts.GetDuration().Seconds() != oldDuration.Seconds() {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
	nnts := nts.SplitByRateInterval(i)
	nsplitTime := time.Date(2012, time.February, 3, 17, 03, 30, 0, time.UTC)
	if nts.TimeStart != splitTime || nts.TimeEnd != nsplitTime {
		t.Error("Incorrect first half", nts)
	}
	if nnts.TimeStart != nsplitTime || nnts.TimeEnd != t2 {
		t.Error("Incorrect second half", nnts)
	}
	if nts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}

	if nts.GetDuration().Seconds() != 150 || nnts.GetDuration().Seconds() != 30 {
		t.Error("Wrong durations.for RateIntervals", nts.GetDuration().Seconds(), nnts.GetDuration().Seconds())
	}
}

func TestTimespanSplitMultipleGroup(t *testing.T) {
	i := &RateInterval{
		EndTime: "17:05:00",
		Rates:   RateGroups{&Rate{0, 2, 1 * time.Second, 1 * time.Second}, &Rate{60 * time.Second, 1, 1 * time.Second, 1 * time.Second}, &Rate{180 * time.Second, 1, 1 * time.Second, 1 * time.Second}},
	}
	t1 := time.Date(2012, time.February, 3, 17, 00, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 3, 17, 04, 0, 0, time.UTC)
	ts := &TimeSpan{TimeStart: t1, TimeEnd: t2, CallDuration: 240 * time.Second}
	oldDuration := ts.GetDuration()
	nts := ts.SplitByRateInterval(i)
	splitTime := time.Date(2012, time.February, 3, 17, 01, 00, 0, time.UTC)
	if ts.TimeStart != t1 || ts.TimeEnd != splitTime {
		t.Error("Incorrect first half", nts)
	}
	if nts.TimeStart != splitTime || nts.TimeEnd != t2 {
		t.Error("Incorrect second half", nts)
	}
	if ts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}

	if ts.GetDuration().Seconds() != 60 || nts.GetDuration().Seconds() != 180 {
		t.Error("Wrong durations.for RateIntervals", ts.GetDuration().Seconds(), nts.GetDuration().Seconds())
	}
	if ts.GetDuration().Seconds()+nts.GetDuration().Seconds() != oldDuration.Seconds() {
		t.Errorf("The duration has changed: %v + %v != %v", ts.GetDuration().Seconds(), nts.GetDuration().Seconds(), oldDuration.Seconds())
	}
	nnts := nts.SplitByRateInterval(i)
	nsplitTime := time.Date(2012, time.February, 3, 17, 03, 00, 0, time.UTC)
	if nts.TimeStart != splitTime || nts.TimeEnd != nsplitTime {
		t.Error("Incorrect first half", nts)
	}
	if nnts.TimeStart != nsplitTime || nnts.TimeEnd != t2 {
		t.Error("Incorrect second half", nnts)
	}
	if nts.RateInterval != i {
		t.Error("RateInterval not attached correctly")
	}

	if nts.GetDuration().Seconds() != 120 || nnts.GetDuration().Seconds() != 60 {
		t.Error("Wrong durations.for RateIntervals", nts.GetDuration().Seconds(), nnts.GetDuration().Seconds())
	}
}

func TestTimespanExpandingPastEnd(t *testing.T) {
	timespans := []*TimeSpan{
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			RateInterval: &RateInterval{Rates: RateGroups{
				&Rate{RateIncrement: 60 * time.Second},
			}},
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 45, 0, time.UTC),
		},
	}
	cd := &CallDescriptor{}
	timespans = cd.roundTimeSpansToIncrement(timespans)
	if len(timespans) != 1 {
		t.Error("Error removing overlaped intervals: ", timespans)
	}
	if !timespans[0].TimeEnd.Equal(time.Date(2013, 9, 10, 14, 31, 0, 0, time.UTC)) {
		t.Error("Error expanding timespan: ", timespans[0])
	}
}

/*
func TestTimespanExpandingCallDuration(t *testing.T) {
	timespans := []*TimeSpan{
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			RateInterval: &RateInterval{Rates: RateGroups{
				&Rate{RateIncrement: 60 * time.Second},
			}},
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 45, 0, time.UTC),
		},
	}
	cd := &CallDescriptor{}
	timespans = cd.roundTimeSpansToIncrement(timespans)

	if timespans[0].CallDuration != time.Minute {
		t.Error("Error setting call duration: ", timespans[0])
	}
}
*/
func TestTimespanExpandingRoundingPastEnd(t *testing.T) {
	timespans := []*TimeSpan{
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 20, 0, time.UTC),
			RateInterval: &RateInterval{Rates: RateGroups{
				&Rate{RateIncrement: 15 * time.Second},
			}},
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 20, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 40, 0, time.UTC),
		},
	}
	cd := &CallDescriptor{}
	timespans = cd.roundTimeSpansToIncrement(timespans)
	if len(timespans) != 2 {
		t.Error("Error removing overlaped intervals: ", timespans)
	}
	if !timespans[0].TimeEnd.Equal(time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC)) {
		t.Error("Error expanding timespan: ", timespans[0])
	}
}

func TestTimespanExpandingPastEndMultiple(t *testing.T) {
	timespans := []*TimeSpan{
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			RateInterval: &RateInterval{Rates: RateGroups{
				&Rate{RateIncrement: 60 * time.Second},
			}},
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 40, 0, time.UTC),
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 40, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 50, 0, time.UTC),
		},
	}
	cd := &CallDescriptor{}
	timespans = cd.roundTimeSpansToIncrement(timespans)
	if len(timespans) != 1 {
		t.Error("Error removing overlaped intervals: ", timespans)
	}
	if !timespans[0].TimeEnd.Equal(time.Date(2013, 9, 10, 14, 31, 0, 0, time.UTC)) {
		t.Error("Error expanding timespan: ", timespans[0])
	}
}

func TestTimespanExpandingPastEndMultipleEqual(t *testing.T) {
	timespans := []*TimeSpan{
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			RateInterval: &RateInterval{Rates: RateGroups{
				&Rate{RateIncrement: 60 * time.Second},
			}},
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 40, 0, time.UTC),
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 40, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 31, 00, 0, time.UTC),
		},
	}
	cd := &CallDescriptor{}
	timespans = cd.roundTimeSpansToIncrement(timespans)
	if len(timespans) != 1 {
		t.Error("Error removing overlaped intervals: ", timespans)
	}
	if !timespans[0].TimeEnd.Equal(time.Date(2013, 9, 10, 14, 31, 0, 0, time.UTC)) {
		t.Error("Error expanding timespan: ", timespans[0])
	}
}

func TestTimespanExpandingBeforeEnd(t *testing.T) {
	timespans := []*TimeSpan{
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			RateInterval: &RateInterval{Rates: RateGroups{
				&Rate{RateIncrement: 45 * time.Second},
			}},
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 31, 0, 0, time.UTC),
		},
	}
	cd := &CallDescriptor{}
	timespans = cd.roundTimeSpansToIncrement(timespans)
	if len(timespans) != 2 {
		t.Error("Error removing overlaped intervals: ", timespans)
	}
	if !timespans[0].TimeEnd.Equal(time.Date(2013, 9, 10, 14, 30, 45, 0, time.UTC)) ||
		!timespans[1].TimeStart.Equal(time.Date(2013, 9, 10, 14, 30, 45, 0, time.UTC)) ||
		!timespans[1].TimeEnd.Equal(time.Date(2013, 9, 10, 14, 31, 0, 0, time.UTC)) {
		t.Error("Error expanding timespan: ", timespans[0])
	}
}

func TestTimespanExpandingBeforeEndMultiple(t *testing.T) {
	timespans := []*TimeSpan{
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			RateInterval: &RateInterval{Rates: RateGroups{
				&Rate{RateIncrement: 45 * time.Second},
			}},
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 30, 50, 0, time.UTC),
		},
		&TimeSpan{
			TimeStart: time.Date(2013, 9, 10, 14, 30, 50, 0, time.UTC),
			TimeEnd:   time.Date(2013, 9, 10, 14, 31, 00, 0, time.UTC),
		},
	}
	cd := &CallDescriptor{}
	timespans = cd.roundTimeSpansToIncrement(timespans)
	if len(timespans) != 3 {
		t.Error("Error removing overlaped intervals: ", timespans)
	}
	if !timespans[0].TimeEnd.Equal(time.Date(2013, 9, 10, 14, 30, 45, 0, time.UTC)) ||
		!timespans[1].TimeStart.Equal(time.Date(2013, 9, 10, 14, 30, 45, 0, time.UTC)) ||
		!timespans[1].TimeEnd.Equal(time.Date(2013, 9, 10, 14, 30, 50, 0, time.UTC)) {
		t.Error("Error expanding timespan: ", timespans[0])
	}
}

func TestTimespanCreateSecondsSlice(t *testing.T) {
	ts := &TimeSpan{
		TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
		TimeEnd:   time.Date(2013, 9, 10, 14, 30, 30, 0, time.UTC),
		RateInterval: &RateInterval{Rates: RateGroups{
			&Rate{Value: 2.0},
		}},
	}
	ts.createIncrementsSlice()
	if len(ts.Increments) != 30 {
		t.Error("Error creating second slice: ", ts.Increments)
	}
	if ts.Increments[0].Cost != 2.0 {
		t.Error("Wrong second slice: ", ts.Increments[0])
	}
}

/*
func TestTimespanCreateSecondsFract(t *testing.T) {
	ts := &TimeSpan{
		TimeStart: time.Date(2013, 9, 10, 14, 30, 0, 0, time.UTC),
		TimeEnd:   time.Date(2013, 9, 10, 14, 30, 30, 100000000, time.UTC),
		RateInterval: &RateInterval{
			RoundingMethod:   utils.ROUNDING_MIDDLE,
			RoundingDecimals: 2,
			Rates: RateGroups{
				&Rate{Value: 2.0},
			},
		},
	}
	ts.createIncrementsSlice()
	if len(ts.Increments) != 31 {
		t.Error("Error creating second slice: ", ts.Increments)
	}
	if len(ts.Increments) < 31 || ts.Increments[30].Cost != 0.2 {
		t.Error("Wrong second slice: ", ts.Increments)
	}
}

func TestTimespanSplitByIncrement(t *testing.T) {
	ts := &TimeSpan{
		TimeStart:    time.Date(2013, 9, 19, 18, 30, 0, 0, time.UTC),
		TimeEnd:      time.Date(2013, 9, 19, 18, 30, 30, 0, time.UTC),
		CallDuration: 50 * time.Second,
	}
	i := &Increment{Duration: time.Second}
	newTs := ts.SplitByIncrement(5, i)
	if ts.GetDuration() != 5*time.Second || newTs.GetDuration() != 25*time.Second {
		t.Error("Error spliting by second: ", ts.GetDuration(), newTs.GetDuration())
	}
	if ts.CallDuration != 25*time.Second || newTs.CallDuration != 50*time.Second {
		t.Error("Error spliting by second at setting call duration: ", ts.GetDuration(), newTs.GetDuration())
	}
}
*/
