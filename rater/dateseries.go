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

package rater

import (
	"sort"
	"strconv"
	"strings"
	"time"
	"fmt"
	"reflect"
)

// Defines years days series
type Years []int

func (ys Years) Sort() {
	sort.Sort(ys)
}

func (ys Years) Len() int {
	return len(ys)
}

func (ys Years) Swap(i, j int) {
	ys[i], ys[j] = ys[j], ys[i]
}

func (ys Years) Less(j, i int) bool {
	return ys[j] < ys[i]
}

// Return true if the specified date is inside the series
func (ys Years) Contains(year int) (result bool) {
	result = false
	for _, yss := range ys {
		if yss == year {
			result = true
			break
		}
	}
	return
}

// Parse Years elements from string separated by sep.
func (ys *Years) Parse(input, sep string) {
	switch input {
	case "*all", "":
		*ys = []int{}
	default:
		elements := strings.Split(input, sep)
		for _, yss := range elements {
			if year, err := strconv.Atoi(yss); err == nil {
				*ys = append(*ys, year)
			}
		}
	}
}

func (ys Years) Serialize( sep string ) string {
	if len(ys) == 0 {
		return "*all"
	}
	var yStr string
	for idx, yr := range ys {
		if idx != 0 {
			yStr = fmt.Sprintf("%s%s%d", yStr, sep, yr)
		} else {
			yStr = strconv.Itoa(yr)
		}
	}
	return yStr
}

var allMonths []time.Month = []time.Month{time.January, time.February, time.March, time.April, time.May, time.June,
			time.July, time.August, time.September, time.October, time.November, time.December}

// Defines months series
type Months []time.Month

func (m Months) Sort() {
	sort.Sort(m)
}

func (m Months) Len() int {
	return len(m)
}

func (m Months) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Months) Less(j, i int) bool {
	return m[j] < m[i]
}

// Return true if the specified date is inside the series
func (m Months) Contains(month time.Month) (result bool) {
	for _, ms := range m {
		if ms == month {
			result = true
			break
		}
	}
	return
}

// Loades Month elemnents from a string separated by sep.
func (m *Months) Parse(input, sep string) {
	switch input {
	case "*all":
		*m = allMonths
	case "*none": // Apier cannot receive empty string, hence using meta-tag
		*m = []time.Month{}
	case "":
		*m = []time.Month{}
	default:
		elements := strings.Split(input, sep)
		for _, ms := range elements {
			if month, err := strconv.Atoi(ms); err == nil {
				*m = append(*m, time.Month(month))
			}
		}
	}
}

// Dumps the months in a serialized string, similar to the one parsed
func (m Months) Serialize( sep string ) string {
	if len(m) == 0 {
		return "*none"
	}
	if reflect.DeepEqual( m, Months(allMonths) ) {
		return "*all"
	}
	var mStr string
	for idx, mt := range m {
		if idx != 0 {
			mStr = fmt.Sprintf("%s%s%d", mStr, sep, mt)
		} else {
			mStr = strconv.Itoa(int(mt))
		}
	}
	return mStr
}


var allMonthDays []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

// Defines month days series
type MonthDays []int

func (md MonthDays) Sort() {
	sort.Sort(md)
}

func (md MonthDays) Len() int {
	return len(md)
}

func (md MonthDays) Swap(i, j int) {
	md[i], md[j] = md[j], md[i]
}

func (md MonthDays) Less(j, i int) bool {
	return md[j] < md[i]
}

// Return true if the specified date is inside the series
func (md MonthDays) Contains(monthDay int) (result bool) {
	result = false
	for _, mds := range md {
		if mds == monthDay {
			result = true
			break
		}
	}
	return
}

// Parse MonthDay elements from string separated by sep.
func (md *MonthDays) Parse(input, sep string) {
	switch input {
	case "*all":
		*md = allMonthDays
	case "":
		*md = []int{}
	default:
		elements := strings.Split(input, sep)
		for _, mds := range elements {
			if day, err := strconv.Atoi(mds); err == nil {
				*md = append(*md, day)
			}
		}
	}
}

// Dumps the month days in a serialized string, similar to the one parsed
func (md MonthDays) Serialize( sep string ) string {
	if len(md) == 0 {
		return "*none"
	}
	if reflect.DeepEqual(md, MonthDays(allMonthDays)) {
		return "*all"
	}
	var mdsStr string
	for idx, mDay := range md {
		if idx != 0 {
			mdsStr = fmt.Sprintf("%s%s%d", mdsStr, sep, mDay)
		} else {
			mdsStr = strconv.Itoa(mDay)
		}
	}
	return mdsStr
}

var allWeekDays []time.Weekday = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
// Defines week days series
type WeekDays []time.Weekday

func (wd WeekDays) Sort() {
	sort.Sort(wd)
}

func (wd WeekDays) Len() int {
	return len(wd)
}

func (wd WeekDays) Swap(i, j int) {
	wd[i], wd[j] = wd[j], wd[i]
}

func (wd WeekDays) Less(j, i int) bool {
	return wd[j] < wd[i]
}

// Return true if the specified date is inside the series
func (wd WeekDays) Contains(weekDay time.Weekday) (result bool) {
	result = false
	for _, wds := range wd {
		if wds == weekDay {
			result = true
			break
		}
	}
	return
}

func (wd *WeekDays) Parse(input, sep string) {
	switch input {
	case "*all":
		*wd = allWeekDays
	case "":
		*wd = []time.Weekday{}
	default:
		elements := strings.Split(input, sep)
		for _, wds := range elements {
			if day, err := strconv.Atoi(wds); err == nil {
				*wd = append(*wd, time.Weekday(day%7)) // %7 for sunday = 7 normalization
			}
		}
	}
}

// Dumps the week days in a serialized string, similar to the one parsed
func (wd WeekDays) Serialize( sep string ) string {
	if len(wd) == 0 {
		return "*none"
	}
	if reflect.DeepEqual( wd, WeekDays(allWeekDays) ) {
		return "*all"
	}
	var wdStr string
	for idx, d := range wd {
		if idx != 0 {
			wdStr = fmt.Sprintf("%s%s%d", wdStr, sep, d)
		} else {
			wdStr = strconv.Itoa(int(d))
		}
	}
	return wdStr
}
