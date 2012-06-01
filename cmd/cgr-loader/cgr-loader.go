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
	"flag"
)

var (
	separator         = flag.String("separator", ",", "Default field separator")
	redisserver       = flag.String("redisserver", "tcp:127.0.0.1:6379", "redis server address (tcp:127.0.0.1:6379)")
	redisdb           = flag.Int("rdb", 10, "redis database number (10)")
	redispass         = flag.String("pass", "", "redis database password")
	monthsFn          = flag.String("month", "Months.csv", "Months file")
	monthdaysFn       = flag.String("monthdays", "MonthDays.csv", "Month days file")
	weekdaysFn        = flag.String("weekdays", "WeekDays.csv", "Week days file")
	destinationsFn    = flag.String("destinations", "Destinations.csv", "Destinations file")
	ratesFn           = flag.String("rates", "Rates.csv", "Rates file")
	ratestimingFn     = flag.String("ratestiming", "RatesTiming.csv", "Rates timing file")
	ratingprofilesFn  = flag.String("ratingprofiles", "RatingProfiles.csv", "Rating profiles file")
	volumediscountsFn = flag.String("volumediscounts", "VolumeDiscounts.csv", "Volume discounts file")
	volumeratesFn     = flag.String("volumerates", "VolumeRates.csv", "Volume rates file")
	inboundbonusesFn  = flag.String("inboundbonuses", "InboundBonuses.csv", "Inbound bonuses file")
	outboundbonusesFn = flag.String("outboundbonuses", "OutboundBonuses.csv", "Outound bonuses file")
	recurrentdebitsFn = flag.String("recurrentdebits", "RecurrentDebits.csv", "Recurent debits file")
	recurrenttopupsFn = flag.String("recurrenttopups", "RecurrentTopups.csv", "Recurent topups file")
	balanceprofilesFn = flag.String("balanceprofiles", "BalanceProfiles.csv", "Balance profiles file")
	sep               rune
)

func main() {
	flag.Parse()
	sep = []rune(*separator)[0]
	loadDataSeries()
	loadDestinations()
	loadRates()
	loadRatesTiming()
	loadRatingProfiles()
	loadVolumeDicounts()
	loadVolumeRates()
	loadInboundBonuses()
	loadOutboundBonuses()
	loadRecurrentDebits()
	loadRecurrentTopups()
	loadBalanceProfiles()
}
