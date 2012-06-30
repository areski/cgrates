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
	// "log"
	"testing"
	"time"
)

func init() {
	sg, _ := NewRedixStorage("tcp:127.0.0.1:6379", 10)
	SetStorageGetter(sg)
}

func TestSingleResultMerge(t *testing.T) {
	t1 := time.Date(2012, time.February, 2, 17, 00, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 2, 17, 01, 0, 0, time.UTC)
	cd := &CallDescriptor{Direction: "OUT", TOR: "0", Tenant: "vdf", Subject: "rif", Destination: "0256", TimeStart: t1, TimeEnd: t2}
	cc1, _ := cd.GetCost()
	if cc1.Cost != 60 {
		t.Errorf("expected 60 was %v", cc1.Cost)
	}
	t1 = time.Date(2012, time.February, 2, 17, 01, 0, 0, time.UTC)
	t2 = time.Date(2012, time.February, 2, 17, 02, 0, 0, time.UTC)
	cd = &CallDescriptor{Direction: "OUT", TOR: "0", Tenant: "vdf", Subject: "rif", Destination: "0256", TimeStart: t1, TimeEnd: t2}
	cc2, _ := cd.GetCost()
	if cc2.Cost != 60 {
		t.Errorf("expected 60 was %v", cc2.Cost)
	}
	cc1.Merge(cc2)
	if len(cc1.Timespans) != 1 || cc1.Timespans[0].GetDuration().Seconds() != 120 {
		t.Error("wrong resulted timespan: ", len(cc1.Timespans), cc1.Timespans[0].GetDuration().Seconds())
	}
	if cc1.Cost != 120 {
		t.Errorf("Exdpected 120 was %v", cc1.Cost)
	}
}

func TestMultipleResultMerge(t *testing.T) {
	t1 := time.Date(2012, time.February, 2, 17, 59, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 2, 18, 00, 0, 0, time.UTC)
	cd := &CallDescriptor{Direction: "OUT", TOR: "0", Tenant: "vdf", Subject: "rif", Destination: "0256", TimeStart: t1, TimeEnd: t2}
	cc1, _ := cd.GetCost()
	if cc1.Cost != 60 {
		t.Errorf("expected 60 was %v", cc1.Cost)
		for _, ts := range cc1.Timespans {
			t.Log(ts.Interval)
		}
	}
	t1 = time.Date(2012, time.February, 2, 18, 00, 0, 0, time.UTC)
	t2 = time.Date(2012, time.February, 2, 18, 01, 0, 0, time.UTC)
	cd = &CallDescriptor{Direction: "OUT", TOR: "0", Tenant: "vdf", Subject: "rif", Destination: "0256", TimeStart: t1, TimeEnd: t2}
	cc2, _ := cd.GetCost()
	if cc2.Cost != 30 {
		t.Errorf("expected 30 was %v", cc2.Cost)
		for _, i := range cd.ActivationPeriods[1].Intervals {
			t.Log(i)
		}
		for _, ts := range cc1.Timespans {
			t.Log(ts.Interval)
		}
	}
	cc1.Merge(cc2)
	if len(cc1.Timespans) != 2 || cc1.Timespans[0].GetDuration().Seconds() != 60 {
		t.Error("wrong resulted timespan: ", len(cc1.Timespans), cc1.Timespans[0].GetDuration().Seconds())
	}
	if cc1.Cost != 90 {
		t.Errorf("Exdpected 90 was %v", cc1.Cost)
	}
}

func TestMultipleInputLeftMerge(t *testing.T) {
	t1 := time.Date(2012, time.February, 2, 17, 59, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 2, 18, 01, 0, 0, time.UTC)
	cd := &CallDescriptor{Direction: "OUT", TOR: "0", Tenant: "vdf", Subject: "rif", Destination: "0256", TimeStart: t1, TimeEnd: t2}
	cc1, _ := cd.GetCost()
	if cc1.Cost != 18 {
		t.Errorf("expected 12 was %v", cc1.Cost)
	}
	t1 = time.Date(2012, time.February, 2, 18, 01, 0, 0, time.UTC)
	t2 = time.Date(2012, time.February, 2, 18, 02, 0, 0, time.UTC)
	cd = &CallDescriptor{Direction: "OUT", TOR: "0", Tenant: "vdf", Subject: "rif", Destination: "0256", TimeStart: t1, TimeEnd: t2}
	cc2, _ := cd.GetCost()
	if cc2.Cost != 6 {
		t.Errorf("expected 6 was %v", cc2.Cost)
	}
	cc1.Merge(cc2)
	if len(cc1.Timespans) != 2 || cc1.Timespans[1].GetDuration().Seconds() != 120 {
		t.Error("wrong resulted timespan: ", len(cc1.Timespans), cc1.Timespans[0].GetDuration().Seconds())
	}
	if cc1.Cost != 24 {
		t.Errorf("Exdpected 24 was %v", cc1.Cost)
	}
}

func TestMultipleInputRightMerge(t *testing.T) {
	t1 := time.Date(2012, time.February, 2, 17, 58, 0, 0, time.UTC)
	t2 := time.Date(2012, time.February, 2, 17, 59, 0, 0, time.UTC)
	cd := &CallDescriptor{Direction: "OUT", TOR: "0", Tenant: "vdf", Subject: "rif", Destination: "0256", TimeStart: t1, TimeEnd: t2}
	cc1, _ := cd.GetCost()
	if cc1.Cost != 12 {
		t.Errorf("expected 12 was %v", cc1.Cost)
	}
	t1 = time.Date(2012, time.February, 2, 17, 59, 0, 0, time.UTC)
	t2 = time.Date(2012, time.February, 2, 18, 01, 0, 0, time.UTC)
	cd = &CallDescriptor{Direction: "OUT", TOR: "0", Tenant: "vdf", Subject: "rif", Destination: "0256", TimeStart: t1, TimeEnd: t2}
	cc2, _ := cd.GetCost()
	if cc2.Cost != 18 {
		t.Errorf("expected 18 was %v", cc2.Cost)
	}
	cc1.Merge(cc2)
	if len(cc1.Timespans) != 2 || cc1.Timespans[0].GetDuration().Seconds() != 120 {
		t.Error("wrong resulted timespan: ", len(cc1.Timespans), cc1.Timespans[0].GetDuration().Seconds())
		t.Log(cc1.Timespans[0].GetDuration())
	}
	if cc1.Cost != 30 {
		t.Errorf("Exdpected 30 was %v", cc1.Cost)
	}
}
