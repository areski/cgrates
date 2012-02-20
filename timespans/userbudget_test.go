package timespans

import (
	"testing"
)

var (
	nationale = &Destination{Id: "nationale", Prefixes: []string{"0257", "0256", "0723"}}
	retea     = &Destination{Id: "retea", Prefixes: []string{"0723", "0724"}}
)

func TestGetSeconds(t *testing.T) {
	b1 := &MinuteBucket{Seconds: 10, Priority: 10, destination: nationale}
	b2 := &MinuteBucket{Seconds: 100, Priority: 20, destination: retea}
	tf1 := &TariffPlan{MinuteBuckets: []*MinuteBucket{b1, b2}}

	ub1 := &UserBudget{id: "rif", minuteBuckets: []*MinuteBucket{b1, b2}, credit: 200, tariffPlan: tf1, resetDayOfTheMonth: 10}
	seconds := ub1.GetSecondsForPrefix(nil, "0723")
	expected := 100
	if seconds != expected {
		t.Errorf("Expected %v was %v", expected, seconds)
	}
}

func TestGetPricedSeconds(t *testing.T) {
	b1 := &MinuteBucket{Seconds: 10, Price: 10, Priority: 10, destination: nationale}
	b2 := &MinuteBucket{Seconds: 100, Price: 1, Priority: 20, destination: retea}
	tf1 := &TariffPlan{MinuteBuckets: []*MinuteBucket{b1, b2}}

	ub1 := &UserBudget{id: "rif", minuteBuckets: []*MinuteBucket{b1, b2}, credit: 21, tariffPlan: tf1, resetDayOfTheMonth: 10}
	seconds := ub1.GetSecondsForPrefix(nil, "0723")
	expected := 21
	if seconds != expected {
		t.Errorf("Expected %v was %v", expected, seconds)
	}
}

/*********************************** Benchmarks *******************************/

func BenchmarkGetSecondForPrefix(b *testing.B) {
	b.StopTimer()
	b1 := &MinuteBucket{Seconds: 10, Price: 10, Priority: 10, destination: nationale}
	b2 := &MinuteBucket{Seconds: 100, Price: 1, Priority: 20, destination: retea}
	tf1 := &TariffPlan{MinuteBuckets: []*MinuteBucket{b1, b2}}

	ub1 := &UserBudget{id: "rif", minuteBuckets: []*MinuteBucket{b1, b2}, credit: 21, tariffPlan: tf1, resetDayOfTheMonth: 10}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ub1.GetSecondsForPrefix(nil,"0723")
	}
}
