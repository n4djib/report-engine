package patch

import "sync/atomic"

// ClientClock is a per-session monotonic counter
// Assign to Patch. ClientSeq BEFORE any other operation
// Call Update() whenever you receive a confirmed server sequence number
type ClientClock struct{ seq uint64 }

// Next return the next sequence number and advances the clock
// Call BEFORE assigning to Patch.ClientSeq
// Never call this more than once per patch - each patch get exactly one value
func (c *ClientClock) Next() uint64 {
	// atomic.AddUint64 is an  atomic read-modify-write operation
	// It is safe to call from multiple goroutines simultaneously
	return atomic.AddUint64(&c.seq, 1)
}

// Update advances the clock past an observed sequence number
// Call this whenever you receive a confirmed sequence number from the server
// this ensures our local clock is alwayse ahead of what the server has seen
func (c *ClientClock) Update(observed uint64) {
	for {
		cur := c.Current()
		if observed <= cur {
			return // already ahead
		}
		// Try to set seq to observed+1 (stay ahead of server)
		if atomic.CompareAndSwapUint64(&c.seq, cur, observed+1) {
			return
		}
	}
}

// Current returns the current value without advancing the clock
// Used for debugging and status display only
func (c *ClientClock) Current() uint64 {
	return atomic.LoadUint64(&c.seq)
}
