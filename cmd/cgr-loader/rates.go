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
package main

import (
	"encoding/csv"
	"github.com/cgrates/cgrates/timespans"
	"log"
	"fmt"
	"os"
	"time"
)

var (
	months         = make(map[string][]time.Month)
	monthdays      = make(map[string][]int)
	weekdays       = make(map[string][]time.Weekday)
	destinations   []*timespans.Destination
	rates          = make(map[string][]*Rate)
	timings        = make(map[string][]*Timing)
	ratesTimings   = make(map[string][]*RateTiming)
	ratingProfiles = make(map[string][]*timespans.CallDescriptor)
)

func loadDataSeries() {
	// MONTHS
	fp, err := os.Open(*monthsFn)
	if err != nil {
		log.Printf("Could not open months file: %v", err)
	} else {
		csvReader := csv.NewReader(fp)
		csvReader.Comma = sep
		for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
			tag := record[0]
			if tag == "Tag" {
				// skip header line
				continue
			}
			for i, m := range record[1:] {
				if m == "1" {
					months[tag] = append(months[tag], time.Month(i+1))
				}
			}
			log.Print(tag, months[tag])
		}
		fp.Close()
	}
	// MONTH DAYS
	fp, err = os.Open(*monthdaysFn)
	if err != nil {
		log.Printf("Could not open month days file: %v", err)
	} else {
		csvReader := csv.NewReader(fp)
		csvReader.Comma = sep
		for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
			tag := record[0]
			if tag == "Tag" {
				// skip header line
				continue
			}
			for i, m := range record[1:] {
				if m == "1" {
					monthdays[tag] = append(monthdays[tag], i+1)
				}
			}
			log.Print(tag, monthdays[tag])
		}
		fp.Close()
	}
	// WEEK DAYS
	fp, err = os.Open(*weekdaysFn)
	if err != nil {
		log.Printf("Could not open week days file: %v", err)
	} else {
		csvReader := csv.NewReader(fp)
		csvReader.Comma = sep
		for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
			tag := record[0]
			if tag == "Tag" {
				// skip header line
				continue
			}
			for i, m := range record[1:] {
				if m == "1" {
					weekdays[tag] = append(weekdays[tag], time.Weekday(((i + 1) % 7)))
				}
			}
			log.Print(tag, weekdays[tag])
		}
		fp.Close()
	}
}

func loadDestinations() {
	fp, err := os.Open(*destinationsFn)
	if err != nil {
		log.Printf("Could not open destinations file: %v", err)
		return
	}
	defer fp.Close()
	csvReader := csv.NewReader(fp)
	csvReader.Comma = sep
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		var dest *timespans.Destination
		for _, d := range destinations {
			if d.Id == tag {
				dest = d
				break
			}
		}
		if dest == nil {
			dest = &timespans.Destination{Id: tag}
			destinations = append(destinations, dest)
		}
		dest.Prefixes = append(dest.Prefixes, record[1:]...)
	}
	log.Print(destinations)
}

func loadRates() {
	fp, err := os.Open(*ratesFn)
	if err != nil {
		log.Printf("Could not open rates timing file: %v", err)
		return
	}
	defer fp.Close()
	csvReader := csv.NewReader(fp)
	csvReader.Comma = sep
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		r, err := NewRate(record[1], record[2], record[3], record[4], record[5])
		if err != nil {
			continue
		}
		rates[tag] = append(rates[tag], r)
		log.Print(tag, rates[tag])
	}
}

func loadTimings() {
	fp, err := os.Open(*timingsFn)
	if err != nil {
		log.Printf("Could not open timings file: %v", err)
		return
	}
	defer fp.Close()
	csvReader := csv.NewReader(fp)
	csvReader.Comma = sep
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}

		t := NewTiming(record[1:]...)
		timings[tag] = append(timings[tag], t)

		log.Print(tag)
		for _, i := range timings[tag] {
			log.Print(i)
		}
	}
}

func loadRatesTimings() {
	fp, err := os.Open(*ratestimingsFn)
	if err != nil {
		log.Printf("Could not open rates timings file: %v", err)
		return
	}
	defer fp.Close()
	csvReader := csv.NewReader(fp)
	csvReader.Comma = sep
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}

		ts, exists := timings[record[2]]
		if !exists {
			log.Printf("Could not get timing for tag %v", record[2])
			continue
		}

		for _, t := range ts {
			rt := NewRateTiming(record[1], t)
			ratesTimings[tag] = append(ratesTimings[tag], rt)
		}
		log.Print(tag)
		for _, i := range ratesTimings[tag] {
			log.Print(i)
		}
	}
}

func loadRatingProfiles() {
	fp, err := os.Open(*ratingprofilesFn)
	if err != nil {
		log.Printf("Could not open destinations rates file: %v", err)
		return
	}
	defer fp.Close()
	csvReader := csv.NewReader(fp)
	csvReader.Comma = sep
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tenant" {
			// skip header line
			continue
		}
		tenant, tor, subject, fallbacksubject := record[0], record[1], record[2], record[3]
		at, err := time.Parse(time.RFC3339, record[5])
		if err != nil {
			log.Printf("Cannot parse activation time from %v", record[5])
			continue
		}
		rts, exists := ratesTimings[record[4]]
		if !exists {
			log.Printf("Could not get rate timing for tag %v", record[4])
			continue
		}
		for _, rt := range rts { // rates timing
			rs, exists := rates[rt.RatesTag]
			if !exists {
				log.Printf("Could not get rates for tag %v", rt.RatesTag)
				continue
			}
			ap := &timespans.ActivationPeriod{
				ActivationTime: at,
			}
			for _, r := range rs { //rates				
				for _, d := range destinations {
					if d.Id == r.DestinationsTag {
						ap.AddInterval(rt.GetInterval(r))
						for _, p := range d.Prefixes { //destinations
							// Search for a CallDescriptor with the same key
							var cd *timespans.CallDescriptor
							for _, c := range ratingProfiles[p] {
								if c.GetKey() == fmt.Sprintf("%s:%s:%s", tenant, subject, p) {
									cd = c
								}
							}
							if cd == nil {
								cd = &timespans.CallDescriptor{
									Tenant:      tenant,
									TOR:         tor,
									Subject:     subject,
									Destination: p,
								}
								ratingProfiles[p] = append(ratingProfiles[p], cd)
							}
							if fallbacksubject != "" {
								// construct a new cd!!!!
							}
							cd.ActivationPeriods = append(cd.ActivationPeriods, ap)
						}
					}
				}
			}
		}
	}
	for dest, cds := range ratingProfiles {
		log.Print(dest)
		for _, cd := range cds {
			log.Print(cd)
		}
	}
}
