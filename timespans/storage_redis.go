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
	"github.com/simonz05/godis"
	"encoding/json"
)

type RedisStorage struct {
	dbNb int
	db   *godis.Client
}

func NewRedisStorage(address string, db int) (*RedisStorage, error) {
	ndb := godis.New(address, db, "")
	return &RedisStorage{db: ndb, dbNb: db}, nil
}

func (rs *RedisStorage) Close() {
	rs.db.Quit()
}

func (rs *RedisStorage) Flush() error {
	return rs.db.Flushdb()
}

func (rs *RedisStorage) GetActivationPeriodsOrFallback(key string) (aps []*ActivationPeriod, fallbackKey string, err error) {
	//rs.db.Select(rs.dbNb)
	elem, err := rs.db.Get(key)
	if err != nil {
		return
	}
	err = json.Unmarshal(elem, &aps)
	if err != nil {
		err = json.Unmarshal(elem, &fallbackKey)
	}
	return
}

func (rs *RedisStorage) SetActivationPeriodsOrFallback(key string, aps []*ActivationPeriod, fallbackKey string) (err error) {
	//.db.Select(rs.dbNb)
	var result []byte
	if len(aps) > 0 {
		result, err = json.Marshal(aps)
	} else {
		result, err = json.Marshal(fallbackKey)
	}
	return rs.db.Set(key, result)
}

func (rs *RedisStorage) GetDestination(key string) (dest *Destination, err error) {
	//rs.db.Select(rs.dbNb + 1)
	if values, err := rs.db.Get(key); err == nil {
		dest = &Destination{Id: key}
		err = json.Unmarshal(values, dest)
	}
	return
}
func (rs *RedisStorage) SetDestination(dest *Destination) (err error) {
	//rs.db.Select(rs.dbNb + 1)
	result, err := json.Marshal(dest)
	return rs.db.Set(dest.Id, result)
}

func (rs *RedisStorage) GetActions(key string) (as []*Action, err error) {
	//rs.db.Select(rs.dbNb + 2)
	if values, err := rs.db.Get(key); err == nil {
		err = json.Unmarshal(values, as)
	}
	return
}

func (rs *RedisStorage) SetActions(key string, as []*Action) (err error) {
	//rs.db.Select(rs.dbNb + 2)
	result, err := json.Marshal(as)
	return rs.db.Set(key, result)
}

func (rs *RedisStorage) GetUserBalance(key string) (ub *UserBalance, err error) {
	//rs.db.Select(rs.dbNb + 3)
	if values, err := rs.db.Get(key); err == nil {
		ub = &UserBalance{Id: key}
		err = json.Unmarshal(values, ub)
	}
	return
}

func (rs *RedisStorage) SetUserBalance(ub *UserBalance) (err error) {
	//rs.db.Select(rs.dbNb + 3)
	result, err := json.Marshal(ub)
	return rs.db.Set(ub.Id, result)
}

func (rs *RedisStorage) GetActionTiming(key string) (at *ActionTiming, err error) {
	if values, err := rs.db.Get(key); err == nil {
		at = &ActionTiming{Id: key}
		err = json.Unmarshal(values, at)
	}
	return
}

func (rs *RedisStorage) SetActionTiming(at *ActionTiming) (err error) {
	result, err := json.Marshal(at)
	return rs.db.Set(at.Id, result)
}
