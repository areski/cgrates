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
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/cgrates/cgrates/utils"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type CSVReader struct {
	sep               rune
	storage           DataStorage
	readerFunc        func(string, rune) (*csv.Reader, *os.File, error)
	actions           map[string][]*Action
	actionsTimings    map[string][]*ActionTiming
	actionsTriggers   map[string][]*ActionTrigger
	accountActions    []*UserBalance
	destinations      []*Destination
	timings           map[string]*Timing
	rates             map[string]*Rate
	destinationRates  map[string][]*DestinationRate
	activationPeriods map[string]*ActivationPeriod
	ratingProfiles    map[string]*RatingProfile
	// file names
	destinationsFn, ratesFn, destinationratesFn, timingsFn, destinationratetimingsFn, ratingprofilesFn,
	actionsFn, actiontimingsFn, actiontriggersFn, accountactionsFn string
}

func NewFileCSVReader(storage DataStorage, sep rune, destinationsFn, timingsFn, ratesFn, destinationratesFn, destinationratetimingsFn, ratingprofilesFn, actionsFn, actiontimingsFn, actiontriggersFn, accountactionsFn string) *CSVReader {
	c := new(CSVReader)
	c.sep = sep
	c.storage = storage
	c.actions = make(map[string][]*Action)
	c.actionsTimings = make(map[string][]*ActionTiming)
	c.actionsTriggers = make(map[string][]*ActionTrigger)
	c.rates = make(map[string]*Rate)
	c.destinationRates = make(map[string][]*DestinationRate)
	c.timings = make(map[string]*Timing)
	c.activationPeriods = make(map[string]*ActivationPeriod)
	c.ratingProfiles = make(map[string]*RatingProfile)
	c.readerFunc = openFileCSVReader
	c.destinationsFn, c.timingsFn, c.ratesFn, c.destinationratesFn, c.destinationratetimingsFn, c.ratingprofilesFn,
		c.actionsFn, c.actiontimingsFn, c.actiontriggersFn, c.accountactionsFn = destinationsFn, timingsFn,
		ratesFn, destinationratesFn, destinationratetimingsFn, ratingprofilesFn, actionsFn, actiontimingsFn, actiontriggersFn, accountactionsFn
	return c
}

func NewStringCSVReader(storage DataStorage, sep rune, destinationsFn, timingsFn, ratesFn, destinationratesFn, destinationratetimingsFn, ratingprofilesFn, actionsFn, actiontimingsFn, actiontriggersFn, accountactionsFn string) *CSVReader {
	c := NewFileCSVReader(storage, sep, destinationsFn, timingsFn, ratesFn, destinationratesFn, destinationratetimingsFn, ratingprofilesFn, actionsFn, actiontimingsFn, actiontriggersFn, accountactionsFn)
	c.readerFunc = openStringCSVReader
	return c
}

func openFileCSVReader(fn string, comma rune) (csvReader *csv.Reader, fp *os.File, err error) {
	fp, err = os.Open(fn)
	if err != nil {
		return
	}
	csvReader = csv.NewReader(fp)
	csvReader.Comma = comma
	csvReader.TrailingComma = true
	return
}

func openStringCSVReader(data string, comma rune) (csvReader *csv.Reader, fp *os.File, err error) {
	csvReader = csv.NewReader(strings.NewReader(data))
	csvReader.Comma = comma
	csvReader.TrailingComma = true
	return
}

func (csvr *CSVReader) WriteToDatabase(flush, verbose bool) (err error) {
	storage := csvr.storage
	if storage == nil {
		return errors.New("No database connection!")
	}
	if flush {
		storage.Flush()
	}
	if verbose {
		log.Print("Destinations")
	}
	for _, d := range csvr.destinations {
		err = storage.SetDestination(d)
		if err != nil {
			return err
		}
		if verbose {
			log.Print(d.Id, " : ", d.Prefixes)
		}
	}
	if verbose {
		log.Print("Rating profiles")
	}
	for _, rp := range csvr.ratingProfiles {
		err = storage.SetRatingProfile(rp)
		if err != nil {
			return err
		}
		if verbose {
			log.Print(rp.Id)
		}
	}
	if verbose {
		log.Print("Action timings")
	}
	for k, ats := range csvr.actionsTimings {
		err = storage.SetActionTimings(k, ats)
		if err != nil {
			return err
		}
		if verbose {
			log.Println(k)
		}
	}
	if verbose {
		log.Print("Actions")
	}
	for k, as := range csvr.actions {
		err = storage.SetActions(k, as)
		if err != nil {
			return err
		}
		if verbose {
			log.Println(k)
		}
	}
	if verbose {
		log.Print("Account actions")
	}
	for _, ub := range csvr.accountActions {
		err = storage.SetUserBalance(ub)
		if err != nil {
			return err
		}
		if verbose {
			log.Println(ub.Id)
		}
	}
	return
}

func (csvr *CSVReader) LoadDestinations() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.destinationsFn, csvr.sep)
	if err != nil {
		log.Print("Could not load destinations file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {

		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		var dest *Destination
		for _, d := range csvr.destinations {
			if d.Id == tag {
				dest = d
				break
			}
		}
		if dest == nil {
			dest = &Destination{Id: tag}
			csvr.destinations = append(csvr.destinations, dest)
		}
		dest.Prefixes = append(dest.Prefixes, record[1])
	}
	return
}

func (csvr *CSVReader) LoadTimings() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.timingsFn, csvr.sep)
	if err != nil {
		log.Print("Could not load timings file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}

		csvr.timings[tag] = NewTiming(record...)
	}
	return
}

func (csvr *CSVReader) LoadRates() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.ratesFn, csvr.sep)
	if err != nil {
		log.Print("Could not load rates file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		var r *Rate
		r, err = NewRate(record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7])
		if err != nil {
			return err
		}
		csvr.rates[tag] = r
	}
	return
}

func (csvr *CSVReader) LoadDestinationRates() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.destinationratesFn, csvr.sep)
	if err != nil {
		log.Print("Could not load rates file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		r, exists := csvr.rates[record[2]]
		if !exists {
			return errors.New(fmt.Sprintf("Could not get rating for tag %v", record[2]))
		}
		dr := &DestinationRate{
			Tag:             tag,
			DestinationsTag: record[1],
			Rate:            r,
		}

		csvr.destinationRates[tag] = append(csvr.destinationRates[tag], dr)
	}
	return
}

func (csvr *CSVReader) LoadDestinationRateTimings() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.destinationratetimingsFn, csvr.sep)
	if err != nil {
		log.Print("Could not load rate timings file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}

		t, exists := csvr.timings[record[2]]
		if !exists {
			return errors.New(fmt.Sprintf("Could not get timing for tag %v", record[2]))
		}

		rt := NewDestinationRateTiming(record[1], t, record[3])
		drs, exists := csvr.destinationRates[record[1]]
		if !exists {
			return errors.New(fmt.Sprintf("Could not find destination rate for tag %v", record[1]))
		}
		for _, dr := range drs {
			if _, exists := csvr.activationPeriods[tag]; !exists {
				csvr.activationPeriods[tag] = &ActivationPeriod{}
			}
			csvr.activationPeriods[tag].AddIntervalIfNotPresent(rt.GetInterval(dr))
		}
	}
	return
}

func (csvr *CSVReader) LoadRatingProfiles() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.ratingprofilesFn, csvr.sep)
	if err != nil {
		log.Print("Could not load rating profiles file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tenant" {
			// skip header line
			continue
		}
		if len(record) != 7 {
			return errors.New(fmt.Sprintf("Malformed rating profile: %v", record))
		}
		tenant, tor, direction, subject, fallbacksubject := record[0], record[1], record[2], record[3], record[4]
		at, err := time.Parse(time.RFC3339, record[6])
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot parse activation time from %v", record[6]))
		}
		key := fmt.Sprintf("%s:%s:%s:%s", direction, tenant, tor, subject)
		rp, ok := csvr.ratingProfiles[key]
		if !ok {
			rp = &RatingProfile{Id: key}
			csvr.ratingProfiles[key] = rp
		}
		for _, d := range csvr.destinations {
			ap, exists := csvr.activationPeriods[record[5]]
			if !exists {
				return errors.New(fmt.Sprintf("Could not load ratinTiming for tag: %v", record[5]))
			}
			newAP := &ActivationPeriod{ActivationTime: at}
			//copy(newAP.Intervals, ap.Intervals)
			newAP.Intervals = append(newAP.Intervals, ap.Intervals...)
			rp.AddActivationPeriodIfNotPresent(d.Id, newAP)
			if fallbacksubject != "" {
				rp.FallbackKey = fmt.Sprintf("%s:%s:%s:%s", direction, tenant, tor, fallbacksubject)
			}
		}
	}
	return
}

func (csvr *CSVReader) LoadActions() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.actionsFn, csvr.sep)
	if err != nil {
		log.Print("Could not load action triggers file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		units, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not parse action units: %v", err))
		}
		unix, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not parse expiration date: %v", err))
		}
		expDate := time.Unix(unix, 0)
		var a *Action
		if record[2] != MINUTES {
			a = &Action{
				ActionType:     record[1],
				BalanceId:      record[2],
				Direction:      record[3],
				Units:          units,
				ExpirationDate: expDate,
			}
		} else {
			value, err := strconv.ParseFloat(record[8], 64)
			if err != nil {
				return errors.New(fmt.Sprintf("Could not parse action price: %v", err))
			}
			minutesWeight, err := strconv.ParseFloat(record[9], 64)
			if err != nil {
				return errors.New(fmt.Sprintf("Could not parse action minutes weight: %v", err))
			}
			weight, err := strconv.ParseFloat(record[9], 64)
			if err != nil {
				return errors.New(fmt.Sprintf("Could not parse action weight: %v", err))
			}
			a = &Action{
				Id:             utils.GenUUID(),
				ActionType:     record[1],
				BalanceId:      record[2],
				Direction:      record[3],
				Weight:         weight,
				ExpirationDate: expDate,
				MinuteBucket: &MinuteBucket{
					Seconds:        units,
					Weight:         minutesWeight,
					Price:          value,
					PriceType:      record[7],
					DestinationId:  record[6],
					ExpirationDate: expDate,
				},
			}
		}
		csvr.actions[tag] = append(csvr.actions[tag], a)
	}
	return
}

func (csvr *CSVReader) LoadActionTimings() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.actiontimingsFn, csvr.sep)
	if err != nil {
		log.Print("Could not load action triggers file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		_, exists := csvr.actions[record[1]]
		if !exists {
			return errors.New(fmt.Sprintf("ActionTiming: Could not load the action for tag: %v", record[1]))
		}
		t, exists := csvr.timings[record[2]]
		if !exists {
			return errors.New(fmt.Sprintf("ActionTiming: Could not load the timing for tag: %v", record[2]))
		}
		weight, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return errors.New(fmt.Sprintf("ActionTiming: Could not parse action timing weight: %v", err))
		}
		at := &ActionTiming{
			Id:     utils.GenUUID(),
			Tag:    record[2],
			Weight: weight,
			Timing: &Interval{
				Months:    t.Months,
				MonthDays: t.MonthDays,
				WeekDays:  t.WeekDays,
				StartTime: t.StartTime,
			},
			ActionsId: record[1],
		}
		csvr.actionsTimings[tag] = append(csvr.actionsTimings[tag], at)
	}
	return
}

func (csvr *CSVReader) LoadActionTriggers() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.actiontriggersFn, csvr.sep)
	if err != nil {
		log.Print("Could not load action triggers file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		value, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not parse action trigger value: %v", err))
		}
		weight, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not parse action trigger weight: %v", err))
		}
		at := &ActionTrigger{
			Id:             utils.GenUUID(),
			BalanceId:      record[1],
			Direction:      record[2],
			ThresholdType:  record[3],
			ThresholdValue: value,
			DestinationId:  record[5],
			ActionsId:      record[6],
			Weight:         weight,
		}
		csvr.actionsTriggers[tag] = append(csvr.actionsTriggers[tag], at)
	}
	return
}

func (csvr *CSVReader) LoadAccountActions() (err error) {
	csvReader, fp, err := csvr.readerFunc(csvr.accountactionsFn, csvr.sep)
	if err != nil {
		log.Print("Could not load account actions file: ", err)
		// allow writing of the other values
		return nil
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		if record[0] == "Tenant" {
			continue
		}
		tag := fmt.Sprintf("%s:%s:%s", record[2], record[0], record[1])
		aTriggers, exists := csvr.actionsTriggers[record[4]]
		if record[4] != "" && !exists {
			// only return error if there was something ther for the tag
			return errors.New(fmt.Sprintf("Could not get action triggers for tag %v", record[4]))
		}
		ub := &UserBalance{
			Type:           UB_TYPE_PREPAID,
			Id:             tag,
			ActionTriggers: aTriggers,
		}
		csvr.accountActions = append(csvr.accountActions, ub)

		aTimings, exists := csvr.actionsTimings[record[3]]
		if !exists {
			log.Printf("Could not get action timing for tag %v", record[3])
			// must not continue here
		}
		for _, at := range aTimings {
			at.UserBalanceIds = append(at.UserBalanceIds, tag)
		}
	}
	return nil
}
