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

package apier

// This file deals with tp_rates management over APIs

import (
	"errors"
	"fmt"
	"time"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

// Creates a new rate within a tariff plan
func (self *ApierV1) SetTPRate(attrs utils.TPRate, reply *string) error {
	if missing := utils.MissingStructFields(&attrs, []string{"TPid", "RateId", "RateSlots"}); len(missing) != 0 {
		return fmt.Errorf("%s:%v", utils.ERR_MANDATORY_IE_MISSING, missing)
	}
	if exists, err := self.StorDb.ExistsTPRate(attrs.TPid, attrs.RateId); err != nil {
		return fmt.Errorf("%s:%s", utils.ERR_SERVER_ERROR, err.Error())
	} else if exists {
		return errors.New(utils.ERR_DUPLICATE)
	}
	rts := make([]*engine.Rate, len(attrs.RateSlots))
	for idx, rtSlot := range attrs.RateSlots {
		var errParse error
		itrvlStrs := []string{rtSlot.RatedUnits, rtSlot.RateIncrements, rtSlot.GroupIntervalStart}
		itrvls := make([]time.Duration, len(itrvlStrs))
		for idxItrvl, itrvlStr := range itrvlStrs {
			if itrvls[idxItrvl], errParse = time.ParseDuration(itrvlStr); errParse != nil {
				return fmt.Errorf("%s:Parsing interval failed:%s", utils.ERR_SERVER_ERROR, errParse.Error())
			}
		}
		rts[idx] = &engine.Rate{attrs.RateId, rtSlot.ConnectFee, rtSlot.Rate, itrvls[0], itrvls[1], itrvls[2],
			rtSlot.RoundingMethod, rtSlot.RoundingDecimals, rtSlot.Weight}
	}
	if err := self.StorDb.SetTPRates(attrs.TPid, map[string][]*engine.Rate{attrs.RateId: rts}); err != nil {
		return fmt.Errorf("%s:%s", utils.ERR_SERVER_ERROR, err.Error())
	}
	*reply = "OK"
	return nil
}

type AttrGetTPRate struct {
	TPid   string // Tariff plan id
	RateId string // Rate id
}

// Queries specific Rate on tariff plan
func (self *ApierV1) GetTPRate(attrs AttrGetTPRate, reply *utils.TPRate) error {
	if missing := utils.MissingStructFields(&attrs, []string{"TPid", "RateId"}); len(missing) != 0 { //Params missing
		return fmt.Errorf("%s:%v", utils.ERR_MANDATORY_IE_MISSING, missing)
	}
	if rt, err := self.StorDb.GetTPRate(attrs.TPid, attrs.RateId); err != nil {
		return fmt.Errorf("%s:%s", utils.ERR_SERVER_ERROR, err.Error())
	} else if rt == nil {
		return errors.New(utils.ERR_NOT_FOUND)
	} else {
		*reply = *rt
	}
	return nil
}

type AttrGetTPRateIds struct {
	TPid string // Tariff plan id
}

// Queries rate identities on specific tariff plan.
func (self *ApierV1) GetTPRateIds(attrs AttrGetTPRateIds, reply *[]string) error {
	if missing := utils.MissingStructFields(&attrs, []string{"TPid"}); len(missing) != 0 { //Params missing
		return fmt.Errorf("%s:%v", utils.ERR_MANDATORY_IE_MISSING, missing)
	}
	if ids, err := self.StorDb.GetTPRateIds(attrs.TPid); err != nil {
		return fmt.Errorf("%s:%s", utils.ERR_SERVER_ERROR, err.Error())
	} else if ids == nil {
		return errors.New(utils.ERR_NOT_FOUND)
	} else {
		*reply = ids
	}
	return nil
}
