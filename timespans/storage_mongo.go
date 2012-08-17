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
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type MongoStorage struct {
	session *mgo.Session
	db      *mgo.Database
}

func NewMongoStorage(host, port, db, user, pass string) (StorageGetter, error) {
	dial := fmt.Sprintf(host)
	if user != "" && pass != "" {
		dial = fmt.Sprintf("%s:%s@%s", user, pass, dial)
	}
	if port != "" {
		dial += ":" + port
	}
	session, err := mgo.Dial(dial)
	if err != nil {
		Logger.Err(fmt.Sprintf("Could not connect to logger database: %v", err))
		return nil, err
	}
	ndb := session.DB(db)
	session.SetMode(mgo.Monotonic, true)
	index := mgo.Index{Key: []string{"key"}, Background: true}
	err = ndb.C("actions").EnsureIndex(index)
	err = ndb.C("actiontimings").EnsureIndex(index)
	index = mgo.Index{Key: []string{"id"}, Background: true}
	err = ndb.C("ratingprofiles").EnsureIndex(index)
	err = ndb.C("destinations").EnsureIndex(index)
	err = ndb.C("userbalances").EnsureIndex(index)

	return &MongoStorage{db: ndb, session: session}, nil
}

func (ms *MongoStorage) Close() {
	ms.session.Close()
}

func (ms *MongoStorage) Flush() (err error) {
	err = ms.db.C("ratingprofiles").DropCollection()
	if err != nil {
		return
	}
	err = ms.db.C("destinations").DropCollection()
	if err != nil {
		return
	}
	err = ms.db.C("actions").DropCollection()
	if err != nil {
		return
	}
	err = ms.db.C("userbalances").DropCollection()
	if err != nil {
		return
	}
	err = ms.db.C("actiontimings").DropCollection()
	if err != nil {
		return
	}
	return nil
}

type AcKeyValue struct {
	Key   string
	Value []*Action
}

type AtKeyValue struct {
	Key   string
	Value []*ActionTiming
}

type LogCostEntry struct {
	Id       string `bson:"_id,omitempty"`
	CallCost *CallCost
}

type LogTimingEntry struct {
	ActionTiming *ActionTiming
	Actions      []*Action
	LogTime      time.Time
}

type LogTriggerEntry struct {
	ubId          string
	ActionTrigger *ActionTrigger
	Actions       []*Action
	LogTime       time.Time
}

func (ms *MongoStorage) GetRatingProfile(key string) (rp *RatingProfile, err error) {
	rp = new(RatingProfile)
	err = ms.db.C("ratingprofiles").Find(bson.M{"_id": key}).One(&rp)
	return
}

func (ms *MongoStorage) SetRatingProfile(rp *RatingProfile) error {
	return ms.db.C("ratingprofiles").Insert(rp)
}

func (ms *MongoStorage) GetDestination(key string) (result *Destination, err error) {
	result = new(Destination)
	err = ms.db.C("destinations").Find(bson.M{"id": key}).One(result)
	if err != nil {
		result = nil
	}
	return
}

func (ms *MongoStorage) SetDestination(dest *Destination) error {
	return ms.db.C("destinations").Insert(dest)
}

func (ms *MongoStorage) GetActions(key string) (as []*Action, err error) {
	result := AcKeyValue{}
	err = ms.db.C("actions").Find(bson.M{"key": key}).One(&result)
	return result.Value, err
}

func (ms *MongoStorage) SetActions(key string, as []*Action) error {
	return ms.db.C("actions").Insert(&AcKeyValue{Key: key, Value: as})
}

func (ms *MongoStorage) GetUserBalance(key string) (result *UserBalance, err error) {
	result = new(UserBalance)
	err = ms.db.C("userbalances").Find(bson.M{"id": key}).One(result)
	return
}

func (ms *MongoStorage) SetUserBalance(ub *UserBalance) error {
	return ms.db.C("userbalances").Insert(ub)
}

func (ms *MongoStorage) GetActionTimings(key string) (ats []*ActionTiming, err error) {
	result := AtKeyValue{}
	err = ms.db.C("actiontimings").Find(bson.M{"key": key}).One(&result)
	return result.Value, err
}

func (ms *MongoStorage) SetActionTimings(key string, ats []*ActionTiming) error {
	return ms.db.C("actiontimings").Insert(&AtKeyValue{key, ats})
}

func (ms *MongoStorage) GetAllActionTimings() (ats map[string][]*ActionTiming, err error) {
	result := AtKeyValue{}
	iter := ms.db.C("actiontimings").Find(nil).Iter()
	ats = make(map[string][]*ActionTiming)
	for iter.Next(&result) {
		ats[result.Key] = result.Value
	}
	return
}

func (ms *MongoStorage) LogCallCost(uuid string, cc *CallCost) error {
	return ms.db.C("cclog").Insert(&LogCostEntry{uuid, cc})
}

func (ms *MongoStorage) GetCallCostLog(uuid string) (cc *CallCost, err error) {
	result := new(LogCostEntry)
	err = ms.db.C("cclog").Find(bson.M{"_id": uuid}).One(result)
	cc = result.CallCost
	return
}

func (ms *MongoStorage) LogActionTrigger(ubId string, at *ActionTrigger, as []*Action) (err error) {
	return ms.db.C("actlog").Insert(&LogTriggerEntry{ubId, at, as, time.Now()})
}

func (ms *MongoStorage) LogActionTiming(at *ActionTiming, as []*Action) (err error) {
	return ms.db.C("actlog").Insert(&LogTimingEntry{at, as, time.Now()})
}
