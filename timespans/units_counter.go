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
	"strconv"
	"strings"
)

// Amount of a trafic of a certain type
type UnitsCounter struct {
	Direction     string
	BalanceId     string
	Units         float64
	MinuteBuckets []*MinuteBucket
}

// Adds the minutes from the received minute bucket to an existing bucket if the destination
// is the same or ads the minutye bucket to the list if none matches.
func (uc *UnitsCounter) addMinuteBucket(newMb *MinuteBucket) {
	if newMb == nil {
		return
	}
	found := false
	for _, mb := range uc.MinuteBuckets {
		if mb.DestinationId == newMb.DestinationId {
			mb.Seconds += newMb.Seconds
			found = true
			break
		}
	}
	if !found {
		uc.MinuteBuckets = append(uc.MinuteBuckets, newMb)
	}
}

/*
Serializes the unit counter for the storage. Used for key-value storages.
*/
func (uc *UnitsCounter) store() (result string) {
	result += uc.Direction + "/"
	result += uc.BalanceId + "/"
	result += strconv.FormatFloat(uc.Units, 'f', -1, 64) + "/"
	for _, mb := range uc.MinuteBuckets {
		result += mb.store() + ","
	}
	result = strings.TrimRight(result, ",")
	return
}

/*
De-serializes the unit counter for the storage. Used for key-value storages.
*/
func (uc *UnitsCounter) restore(input string) {
	elements := strings.Split(input, "/")
	if len(elements) != 4 {
		return
	}
	uc.Direction = elements[0]
	uc.BalanceId = elements[1]
	uc.Units, _ = strconv.ParseFloat(elements[2], 64)
	for _, mbs := range strings.Split(elements[3], ",") {
		mb := &MinuteBucket{}
		mb.restore(mbs)
		uc.MinuteBuckets = append(uc.MinuteBuckets, mb)
	}
}
