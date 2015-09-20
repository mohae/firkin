// Queue implements a queue that can either be bounded or unbounded.
// Allocations are minimized and the queue is concurrency safe.
//
// On queue creation, the initial capacity will be equal to the received size.
// If the 'bounded' bool is true, it will be a bounded queue, otherwise it can
// grow as needed.
//
// After dequeuing an item, the head position will be checked. If the queue
// is empty, head > tail, the queue will be reset.
//
// Once a queue grows, it will not shrink.
//
// If a bounded queue is at its capacity, a check will be done to see if there
// is any space at the front of the queue: if items have been dequeued. If
// there are, the items will be shifted to make room for the new item, which
// will then be enqueued. If the bounded queue is full, an error will occur.
package dq

import (
	"fmt"
	"sync"
)

// shiftPercent is the default value for shifting the queue items to the
// front of the queue instead of growing the queue. If at least the % of
// the items have been removed from the queue, the items in the queue will
// be shifted to make room; otherwise the queue will grow
var shiftPercent = 50

// Queue represents a queue and everything needed to manage it. The preferred
// method for creating a new Queue is to use either NewQ() or NewQueue().
type Queue struct {
	sync.RWMutex
	items        []interface{}
	head         int // current item in queue
	bounded      bool // whether or not this is a bounded queue
	shiftPercent int // the % of items that need to be removed before shifting occurs
}

// NewQ returns an empty queue with a capacity equal to the recieved size value. If
// maxCap is > 0, the queue will not grow larger than maxCap; if it is at maxCap
// and growth is requred to enqueue an item, an error will occur.
func NewQ(size int, bounded bool) *Queue {
	return &Queue{items: make([]interface{}, 0, size), bounded: bounded, shiftPercent: shiftPercent}
}

// NewQueue is a convenience wrapper to NewQ().
func NewQueue(size int, bounded bool) *Queue {
	return NewQ(size, bounded)
}

// SetShiftPercent sets the queue's shiftPercent: the percentage of the queue
// that must be empty before the remaining items will be shifted to the
// the beginning of the slice. This occurs when the slice is set to grow.
//
// Valid range of values are 0-100, inclusive. Vaues < 0 are set to 0 and
// values > 100 are set to 100.
func (q *Queue) SetShiftPercent(i int) {
		if i < 0 {
			q.shiftPercent = 0
			return
		}
		if i > 100 {
			q.shiftPercent = 100
			return
		}
		q.shiftPercent = i
}

// Enqueue: adds an item to the queue. If adding the item requires growing
// the queue, the queue will either be shifted, to make room at the end of
// the queue, or it will grow.
//
// If the queue is a bounded queue and is full, an error will be returned.
func (q *Queue) Enqueue(item interface{}) error {
	q.Lock()
	defer q.Unlock()
	// See if it needs to grow
	if len(q.items) == cap(q.items) {
		shifted := q.shift()
		// if we weren't able to make room by shifting, grow the queue
		if !shifted {
			if q.bounded {
				return fmt.Errorf("bounded queue full: cannot enqueue '%v'", item)
			}
		}
	}
	q.items = append(q.items, item)
	return nil
}

// Dequeue removes an item from the queue. If the removal of the item empties
// the queue, the head and tail will be set to 0.
func (q *Queue) Dequeue() interface{} {
	q.Lock()
	i := q.items[q.head]
	if q.head == len(q.items) {
		q.Unlock()
		q.Reset()
		return i
	}
	q.head++
	q.Unlock()
	return i
}

// IsEmpty returns whether or not the queue is empty
func (q *Queue) IsEmpty() bool {
	q.RLock()
	if len(q.items) == 0  {
		q.RUnlock()
		return true
	}
	q.RUnlock()
	return false
}

// IsFull returns whether or not the queue is full. This only applies to
// bounded queues
func (q *Queue) IsFull() bool {
	if !q.bounded {
		return false
	}
	if q.head > 0 {
		return false
	}
	if len(q.items) == cap(q.items) {
		return true
	}
	return false
}

// Count returns the current number of items in the queue
func (q *Queue) Count() int {
	q.RLock()
	defer q.RUnlock()
	return len(q.items) - q.head
}

// shift: if either shiftPercent items have been removed from the queue or the
// queue is a boudned queue, the remaining items in the queue will be shifted
// to the beginning of the queue. Returns whether or not a shift occurred
func (q *Queue) shift() bool {
	// shift percent applies to unbounded queues.
	if !q.bounded {
		if q.head < (cap(q.items)*q.shiftPercent)/100 {
			return false
		}
	}
	// if this is a capped queue, but there is no room for a shift, nothing to do.
	if q.bounded {
		if q.head == 0 {
			return false
		}
	}
	q.items = append(q.items[:0], q.items[q.head:]...)
	// set the pointers to the correct position
	q.head = 0
	return true
}

// Reset resets the queue; head and tail point to element 0.
func (q *Queue) Reset() {
	q.Lock()
	q.head = 0
	q.items = q.items[:0]
	q.Unlock()
}
