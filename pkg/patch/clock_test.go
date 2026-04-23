package patch

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestClockMonotonicallyIncreases(t *testing.T) {
	clock := &ClientClock{
		seq: 1,
	}
	val := clock.Next()
	if val != 2 {
		t.Errorf("Expected 2 got %v", val)
	}
	val2 := clock.Next()
	if val2 != 3 {
		t.Errorf("Expected 2 got %v", val2)
	}
}

func TestClockUpdateAdvancesPastObserved(t *testing.T) {
	clock := ClientClock{seq: 5}

	// server says i saw 10
	clock.Update(10)

	final := clock.Next()
	if final != 12 {
		t.Errorf("After update to 10 , and after Next() it should return 12 got %v", final)
	}
}

func TestClockUpdateIgnoresPastValues(t *testing.T) {
	// time of the client
	clock := ClientClock{seq: 50}

	// Server saw only 20
	clock.Update(20)

	// the clock should stay at 50 so Next() should return 51
	final := clock.Next()

	if final != 51 {
		t.Errorf("The Clock should ignored the old update, Expected %v got %v", 51, final)
	}
}

func TestConcurrentClockIncrements(t *testing.T) {
	clock := ClientClock{seq: 0}

	var wg sync.WaitGroup
	numGoroutines := 100
	incrementPerRouting := 1000

	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range incrementPerRouting {
				clock.Next()
			}
		}()
	}
	wg.Wait()

	final := atomic.LoadUint64(&clock.seq)
	if final != uint64(numGoroutines)*uint64(incrementPerRouting) {
		t.Errorf("Expected %v got %v", numGoroutines*incrementPerRouting, final)
	}
}
