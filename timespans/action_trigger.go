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
	"log"
	"sort"
)

type ActionTrigger struct {
	BalanceId      string
	ThresholdValue float64
	DestinationId  string
	Weight         float64
	ActionsId      string
	executed       bool
}

func (at *ActionTrigger) Execute(ub *UserBalance) (err error) {
	userBalancesRWMutex.Lock()
	defer userBalancesRWMutex.Unlock()
	var aac ActionPriotityList
	aac, err = storageGetter.GetActions(at.ActionsId)
	aac.Sort()
	if err != nil {
		log.Print("Failed to get actions: ", err)
		return
	}
	for _, a := range aac {
		actionFunction, exists := actionTypeFuncMap[a.ActionType]
		if !exists {
			log.Printf("Function type %v not available, aborting execution!", a.ActionType)
			return
		}
		err = actionFunction(ub, a)
	}
	at.executed = true
	storageGetter.SetUserBalance(ub)
	return
}

// Structure to store actions according to weight
type ActionTriggerPriotityList []*ActionTrigger

func (atpl ActionTriggerPriotityList) Len() int {
	return len(atpl)
}

func (atpl ActionTriggerPriotityList) Swap(i, j int) {
	atpl[i], atpl[j] = atpl[j], atpl[i]
}

func (atpl ActionTriggerPriotityList) Less(i, j int) bool {
	return atpl[i].Weight < atpl[j].Weight
}

func (atpl ActionTriggerPriotityList) Sort() {
	sort.Sort(atpl)
}