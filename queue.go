// Queue implements a dynamic queue that will grow as needed, minimize
// unnecessary growth, and can be safely accessed from multiple routines.
//
// The queue can be configured with a max capacity, which will cap the size of
// the queue so unbounded grouwth cannot occur.  When a maximum capacity for
// a queue is defined, an error will occur if the queue is full and an attempt
// is made to add another item to the queue.
//
// A queue is created with a minimum length and an optional maximum size
// (capacity).  If the max size of the queue == 0, the queue will be unbounded.
// The growth rate of the queue is similar to that of a slice. When a queue
// grows, all items in the queue are shifted so that the head of the queue
// points to the first element in the queue.
//
// TODO: rewrite the following for clarity
// Before an item is enqueued, the queue is checked to see if thie new item
// will cause it to grow. If the tail == length, growth may occur. If the
// head of the queue is past a certain point in the queue, which is currently
// calculated using a percentage, the items in the queue will be shifted to start at
// the beginning of the slice, instead of growing the slice. The queue's head and tail
// will then be updated to reflect the shift.
//
// After dequeuing an item, the head position will be checked. If the queue
// io empty, head > tail, head and tail will be set to 0. This allows for
// efficient reuse of the queue without having to check to see if the queue
// items should be shifted or the queue should be grown. The contents of the
// queue are not zeroed out.
//
// Once a queue grows, it will not be shrunk. This behavior may change in the
// future.
//
// All exported methods on the queue use locking so that the queue is safe for
// concurrency.  Unexported methods do not do any locking/unlocking since it
// is expected that the calling method has already obtained the lock and will
// release it as appropriate.
package queue

import (
	"fmt"
	"sync"
)

// shiftPercent is the default value for shifting the queue items to the
// front of the queue instead of growing the queue. If at least the % of
// the items have been removed from the queue, the items in the queue will
// be shifted to make room; otherwise the queue will grow
var shiftPercent = 50

// Queue represents a queue and everything needed to manage it. The preferred method
// for creating a new Queue is to use the New() func.
type Queue struct {
	sync.RWMutex
	items        []interface{}
	head         int // current item in queue
	tail         int // tail is the next insert point. last item is tail - 1
	maxCap       int // if > 0, the queue's cap cannot grow beyond this value
	shiftPercent int // the % of items that need to be removed before shifting occurs
}

// New returns an empty queue with a capacity equal to the recieved size value. If
// maxCap is > 0, the queue will not grow larger than maxCap; if it is at maxCap
// and growth is requred to enqueue an item, an error will occur.
func New(size, maxCap int) *Queue {
	return &Queue{items: make([]interface{}, size, size), maxCap: maxCap, shiftPercent: shiftPercent}
}

// Enqueue: adds an item to the queue. If adding the item requires growing
// the queue, the queue will either be shifted, to make room at the end of the queue
// or it will grow. If the queue cannot be grown, an error will be returned.
func (q *Queue) Enqueue(item interface{}) error {
	q.Lock()
	defer q.Unlock()
	// See if it needs to grow
	if q.tail == cap(q.items) {
		shifted := q.shift()
		// if we weren't able to make room by shifting, grow the queue/
		if !shifted {
			err := q.grow()
			if err != nil {
				return err
			}
		}
	}
	if q.tail > len(q.items)  - 1 {
		q.items = append(q.items, item)
	} else {
		q.items[q.tail] = item
	}
	q.tail++
	return nil
}

// Dequeue removes an item from the queue. If the removal of the item empties
// the queue, the head and tail will be set to 0.
func (q *Queue) Dequeue() interface{} {
	q.Lock()
	i := q.items[q.head]
	q.head++
	if q.head > q.tail {
		q.Unlock()
		q.Reset()
		return i
	}
	q.Unlock()
	return i
}

// IsEmpty returns whether or not the queue is empty
func (q *Queue) IsEmpty() bool {
	q.RLock()
	if q.tail == 0 || q.head == q.tail {
		q.RUnlock()
		return true
	}
	q.RUnlock()
	return false
}

// Tail returns the current tail position
func (q *Queue) Tail() int {
	q.RLock()
	defer q.RUnlock()
	return q.tail
}

// Head returns the current head position
func (q *Queue) Head() int {
	q.RLock()
	defer	q.RUnlock()
	return q.head
}

// Length returns the current length(cap) of the queue. Note, this is not the
// number of items in the queue, for that use Items()
func (q *Queue) Length() int {
	q.RLock()
	defer q.RUnlock()
	return cap(q.items)
}

// ItemCount returns the current number of items in the queue
func (q *Queue) ItemCount() int {
	q.RLock()
	defer q.RUnlock()
	return q.tail - q.head
}

// shift: if shiftPercent items have been removed from the queue, the remaining
// items in the queue will be shifted to element 0-n, where n is the number of
// remaining items in the queue. Returns whether or not a shift occurred
func (q *Queue) shift() bool {
	if q.head < (cap(q.items)*q.shiftPercent)/100 {
		return false
	}
	q.items = append(q.items[:0], q.items[q.head:q.tail]...)
	// set the pointers to the correct position
	q.tail = q.tail - q.head
	q.head = 0
	return true
}

// grow grows the slice using an algorithm similar to growSlice(). This is a bit slower
// than relying on slice's automatic growth, but allows for capacity enforcement w/o
// growing the slice cap beyond the configured maxCap, if applicable.
//
// Since a temporary slice is created to store the current queue, all items in queue
// are automatically shifted
func (q *Queue) grow() error {
	if cap(q.items) == q.maxCap && q.maxCap > 0 {
		return fmt.Errorf("groweQueue: cannot grow beyond max capacity of %d", q.maxCap)
	}
	var len int
	if cap(q.items) < 1024 {
		len = cap(q.items) << 1
	} else {
		len = cap(q.items) + cap(q.items)/4
	}
	// If the maxCap is set, cannot grow it beyond that
	if len > q.maxCap && q.maxCap > 0 {
		len = q.maxCap
	}
	// create a new slice of len
	tmp := make([]interface{}, len, len)
	// copy the remaining elements
	copy(tmp, q.items[q.head:q.tail])
	q.items = tmp
	q.tail = q.tail - q.head
	q.head = 0
	return nil
}

// Reset resets the queue; head and tail point to element 0.
func (q *Queue) Reset() {
	q.Lock()
	q.head = 0
	q.tail = 0
	q.Unlock()
}
