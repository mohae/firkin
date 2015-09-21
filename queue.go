// There are two types of queues, a bounded queue (bqueue) and an unbounded
// queue (queue). Both queues are thread-safe.
//
// For bounded queues, an allocation only occurs at queue creation.
//
// For unbounded queues, the initial capacity of the queue will be equal to the
// received size.
package dq

import (
	"sync"
)

// shiftPercent is the default value for shifting the queue items to the
// front of the queue instead of growing the queue.  If at least the % of
// the items have been removed from the queue, the items in the queue will
// be shifted to make room; otherwise the queue will grow.  This only applies
// to unbounded queues and can be set per queue.
var shiftPercent = 50

// Queue represents an unbounded queue and everything needed to manage it.
// The preferred method for creating a new Queue is to use either NewQ()
// or its alias, NewQueue().
type Queue struct {
	sync.Mutex
	items        []interface{}
	head         int // current item in queue
	shiftPercent int // the % of items that need to be removed before shifting occurs
}

// NewQ returns an empty queue with an initial  capacity equal to the recieved
// size.
func NewQ(size int) *Queue {
	return &Queue{items: make([]interface{}, 0, size), shiftPercent: shiftPercent}
}

// NewQueue is a convenience wrapper to NewQ().
func NewQueue(size int) *Queue {
	return NewQ(size)
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
		_ = q.shift()
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
		q.reset()
		return i
	}
	q.head++
	q.Unlock()
	return i
}

// IsEmpty returns whether or not the queue is empty
func (q *Queue) IsEmpty() bool {
	q.Lock()
	if len(q.items) == 0  {
		q.Unlock()
		return true
	}
	q.Unlock()
	return false
}

// Len returns the current number of items in the queue
func (q *Queue) Len() int {
	q.Lock()
	defer q.Unlock()
	return len(q.items) - q.head
}

// Cap returns the current size of the queue
func (q *Queue) Cap() int {
	q.Lock()
	defer q.Unlock()
	return cap(q.items)
}

// shift: if either shiftPercent items have been removed from the queue or the
// queue is a boudned queue, the remaining items in the queue will be shifted
// to the beginning of the queue. Returns whether or not a shift occurred
func (q *Queue) shift() bool {
	if q.head < (cap(q.items)*q.shiftPercent)/100 {
		return false
	}
	q.items = append(q.items[:0], q.items[q.head:]...)
	// set the pointers to the correct position
	q.head = 0
	return true
}

// Reset resets the queue; head and tail point to element 0.
func (q *Queue) reset() {
	q.Lock()
	q.head = 0
	q.items = q.items[:0]
	q.Unlock()
}
