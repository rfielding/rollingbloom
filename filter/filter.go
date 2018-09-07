package filter

import (
	"github.com/willf/bloom"
)

/*
  This is a triple-buffered bloom filter, so that
  - All add and test are done with a logical clock
  - There is a rollover interval to the next filter
  - Current writes are done to this, and next filter
  - When we use a timestamp that rolls us over, the current filter is wiped first.

  Example:

  Suppose that the rollover interval is 60, and with timestamps of ints:

  00 0:  add  "192.168.2.3" (written to 0 and 1)
  03 0:  add  "192.168.2.4"
  10 0:  add  "192.168.2.3"
  12 0:  add  "192.168.5.55"
  18 0:  test "192.168.5.55" (true)
  88 1:  test "192.168.5.55" (written to 1 and 2)
  ..
  122 2: test "192.168.2.3" (false - presuming it was last written at 10)

  Also, we add in a time dependent bit into the hash so that once a wrong result is found,
  it is infeasible to keep using it.
*/

type Rolling struct {
	Filters      []*bloom.BloomFilter
	Current      int
	Period       int64
	M            uint
	K            uint
	CurrentRound int
}

// Allow bloom filter to roll forward to age out
// items, by writing into current and next, while
// resetting prev bloom filter when we cross over into
// next time interval.
func NewRolling(m uint, k uint, period int64) *Rolling {
	f := &Rolling{}
	f.Period = period
	f.M = m
	f.K = k
	f.Filters = make([]*bloom.BloomFilter, 3)
	f.Filters[0] = bloom.New(m, k)
	f.Filters[1] = bloom.New(m, k)
	f.Filters[2] = bloom.New(m, k)
	return f
}

func (f *Rolling) advance(t int64) {
	tcurrentRound := int(t / f.Period)
	tcurrent := tcurrentRound % 3
	if f.CurrentRound != tcurrentRound {
		f.Filters[(tcurrent + 3 - 1)%3] = bloom.New(f.M, f.K)
		f.Current = tcurrent
		f.CurrentRound = tcurrentRound
	}
}

func (f *Rolling) Add(v []byte, t int64) {
	f.advance(t)
	//v1 := append(v, byte(f.CurrentRound%256))
	//v2 := append(v, byte((f.CurrentRound+1)%256))
	f.Filters[f.Current].Add(v)
	f.Filters[(f.Current+1)%3].Add(v)
}

func (f *Rolling) Test(v []byte, t int64) bool {
	f.advance(t)
	//v1 := append(v, byte(f.CurrentRound%256))
	return f.Filters[f.Current].Test(v)
}
