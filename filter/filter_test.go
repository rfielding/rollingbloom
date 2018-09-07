package filter_test

import (
	"math/rand"
	"fmt"
	"testing"
	"github.com/rfielding/rollingbloom/filter"
)

func TestRollover(t *testing.T) {
	// Add items in the filter with a time
	f := filter.NewRolling(20*1000, 5, 60)
	f.Add([]byte("fark"), 25)
	if !f.Test([]byte("fark"), 30) {
		t.Fail()
	}
	f.Add([]byte("welp"), 55)
	f.Add([]byte("heck"), 85)
	if !f.Test([]byte("fark"), 90) {
		t.Fail()
	}
	f.Add([]byte("gaah"), 96)
	f.Add([]byte("blargh"), 150)
	if f.Test([]byte("fark"), 160) {
		t.Fail()
	}
	if !f.Test([]byte("blargh"), 150+60) {
		t.Fail()
	}
	if f.Test([]byte("blargh"), 150+120) {
		t.Fail()
	}
}

func BenchmarkRollover(b *testing.B) {
	f := filter.NewRolling(20*1000, 5, 100)
	t := int64(0)
	coincidences := 0
	added := int64(100000)
	for t < added  {
		if rand.Intn(10) > 5 {
			y := (11*t*t*t + 13*t*t + 7) % 5
			x := ((t*97+13) % 3)
			stmt := fmt.Sprintf("%d = f[%d]", y, x)
			f.Add([]byte(stmt), t)
		} else {
			y := (19*t*t*t + 17*t*t + 11) % 5
			x := (t % 3)
			stmt := fmt.Sprintf("%d = f[%d]", y, x)
			if f.Test([]byte(stmt), t+added) {
				coincidences++
			}
		}
		t++	
	}
	fmt.Printf("coincidences %d/%d: %f\n", coincidences,added,(1.0*float32(coincidences))/float32(added))
}
