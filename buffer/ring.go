// Package buffer provides a thread-safe buffer implementation for a ring
// buffer.
package buffer

import (
	"math"

	"github.com/mohae/firkin/queue"
)

// Ring is a ring buffer implementation wrapping queue.Circular.
type Ring struct {
	queue.Circular
}

// NewRing returns a ring buffer initalized with 'size' slots.
func NewRing(size int) *Ring {
	return &Ring{*queue.NewCircular(size)}
}

// Enqueue enques an item, If the buffer is full, the oldest item will
// be evicted.
func (r *Ring) Enqueue(item interface{}) error {
	r.Lock()
	// if the buffer is full, move the head forward
	if r.isFull() {
		r.Head = int(math.Mod(float64(r.Head+1), float64(cap(r.Items))))
	}
	r.Items[r.Tail] = item
	r.Tail = int(math.Mod(float64(r.Tail+1), float64(cap(r.Items))))
	r.Unlock()
	return nil
}

// isFull is an unexported version that expects the caller to handle locking.
// This eliminates double locking on enqueue
func (r *Ring) isFull() bool {
	if r.Head != int(math.Mod(float64(r.Tail+1), float64(cap(r.Items)))) {
		return false
	}
	return true
}
