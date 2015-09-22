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

// Queuer interface
type Queuer interface {
	Enqueue(item interface{}) error
	Dequeue() (interface{}, bool)
	Peek() (interface{}, bool)
	IsEmpty() bool
	IsFull() bool
	Reset()
}

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
	initCap      int
	items        []interface{}
	head         int // current item in queue
	shiftPercent int // the % of items that need to be removed before shifting occurs
}

// NewQ is a convenience wrapper to NewQ().
func NewQ(size int) *Queue {
	return NewQueue(size)
	}

// NewQueue returns an empty queue with an initial  capacity equal to the
// recieved size.
func NewQueue(size int) *Queue {
	return &Queue{initCap: size, items: make([]interface{}, 0, size), shiftPercent: shiftPercent}
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
// the queue, the head and tail will be set to 0. If the queue is empty, a
// false will be returned, else true.
func (q *Queue) Dequeue() (interface{}, bool) {
	q.Lock()
	defer q.Unlock()
	if q.isEmpty() {
		return nil, false
	}
	q.head++
	return q.items[q.head-1], true
}

// Peek returns the next item in the queue. Post-peek, the queue remains the
// same.
func (q *Queue) Peek() (interface{}, bool) {
	q.Lock()
	defer q.Unlock()
	if q.isEmpty() {
		return nil, false
	}
	return q.items[q.head], true
}

// IsEmpty returns whether or not the queue is empty
func (q *Queue) IsEmpty() bool {
	q.Lock()
	defer q.Unlock()
	return q.isEmpty()
}

// isEmpty is an unexported version that doesn't lock because the caller
// will have handled that. Reduces multiple locks/unlocks during operations
// that need to check for emptiness and have already obtained a lock
func (q *Queue) isEmpty() bool {
	if q.head == len(q.items) {
		return true
	}
	return false
}
// IsFull returns false; this is implemented to fulfill Queuer but a dynamic
// queue will never be full.
func (q *Queue) IsFull() bool {
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

// reset resets the queue; head and tail point to element 0. This does not
// shrink the queue; for that use Reset()
func (q *Queue) reset() {
	q.Lock()
	q.head = 0
	q.items = q.items[:0]
	q.Unlock()
}

// Reset resets the queue to its original capacity.
