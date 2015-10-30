package queue

import (
	"fmt"
	"math"
)

// Circular is a bounded queue implemented as a circular queue.  Even though
// Items, Head, and Tail are exported, in most cases, they should not be
// directly.  Doing so may lead to outcomes less than desirable. Use the
// exported methods to interact with the Circular queue.
type Circular struct {
	Queue
	Tail int
}

// NewCircular returns an initialized circular queue. Even though creating
// the slice with an initial length is much slower than creating one without
// the initial length, cap only, this is done to simplify the actual queue
// management. Don't need to worry about appending vs adding via index and
// don't need to check to see if an append will cause the slice to grow.
//
// The slice is 1 slot larger than the requested size for empty/full
// detection.
func NewCircular(size int) *Circular {
	size++
	c := Circular{Queue: *NewQueue(size)}
	_ = c.zeroQueue()
	return &c
}

// Enqueue will return an error if the queue is full
func (c *Circular) Enqueue(item interface{}) error {
	c.Lock()
	if c.isFull() {
		c.Unlock()
		return fmt.Errorf("queue full: cannot enqueue %v", item)
	}
	c.Items[c.Tail] = item
	c.Tail = int(math.Mod(float64(c.Tail+1), float64(cap(c.Items))))
	c.Unlock()
	return nil
}

// Dequeue will remove an item from the queue and return it. If the queue is
// empty, a false will be returned.
func (c *Circular) Dequeue() (interface{}, bool) {
	c.Lock()
	item, ok := c.peek()
	if ok {
		c.Head = int(math.Mod(float64(c.Head+1), float64(cap(c.Items))))
	}
	c.Unlock()
	return item, ok
}

// Peek will return the next item in the queue without removing it from the
// queue. If the queue is empty, a false will be returned.
func (c *Circular) Peek() (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	return c.peek()
}

// peek is an unexported version that expects the caller to handle locking.
func (c *Circular) peek() (interface{}, bool) {
	if c.isEmpty() {
		return nil, false
	}
	return c.Items[c.Head], true
}

// IsEmpty returns whether or not the queue is empty
func (c *Circular) IsEmpty() bool {
	c.Lock()
	defer c.Unlock()
	return c.isEmpty()
}

// isEmpty is an unexported version that expects the caller to handle locking.
// This eliminates double locking on dequeue and peek
func (c *Circular) isEmpty() bool {
	if c.Head == c.Tail {
		return true
	}
	return false
}

// IsFull returns whether or not the queue is full
func (c *Circular) IsFull() bool {
	c.Lock()
	defer c.Unlock()
	return c.isFull()
}

// isFull is an unexported version that expects the caller to handle locking.
// This eliminates double locking on enqueue
func (c *Circular) isFull() bool {
	if c.Head == int(math.Mod(float64(c.Tail+1), float64(cap(c.Items)))) {
		return true
	}
	return false
}

// Len returns the current length of the queue (# of items in queue)
func (c *Circular) Len() int {
	c.Lock()
	defer c.Unlock()
	return c.plen()
}

// plen returns the current length of the queue (# items in queue).  This
// unexported method does not do any locking of its own, it relies on
// the caller to take care of locking.
func (c *Circular) plen() int {
	l := c.Tail
	if c.Tail < c.Head {
		l += cap(c.Items)
	}
	return l - c.Head
}

// Cap returns the current queue capacity:
//    queue cap = cap(queue) - 1
func (c *Circular) Cap() int {
	c.Lock()
	defer c.Unlock()
	return cap(c.Items) - 1
}

// Resize resizes a queue; zeroing out the slots.
func (c *Circular) Resize(size int) int {
	c.Lock()
	if ((size + 1) == c.InitCap) || (size == 0 && cap(c.Items) == c.InitCap) {
		c.Unlock()
		return size
	}
	// tmp slice of remaining Items
	tmp := make([]interface{}, 0, c.plen())
	for i := 0; i < cap(tmp); i++ {
		tmp = append(tmp, c.Items[c.Head:c.Head])
		c.Head = int(math.Mod(float64(c.Head+1), float64(cap(c.Items))))
	}
	c.Items = tmp
	c.Head = 0
	c.Tail = len(tmp)
	c.Unlock()
	x := c.Queue.Resize(size + 1)
	c.Lock()
	_ = c.zeroQueue()
	c.Unlock()
	return x
}

// Reset resets a queue, zeroing out the remaining slots.
func (c *Circular) Reset() {
	c.Queue.Reset()
	c.Lock()
	c.Tail = 0
	_ = c.zeroQueue()
	c.Unlock()

}

// zeroQueue appends the zero value to the queue unti the queue is at cap.
// This is needed because then length of the after a queue.Resize() or
// queue.Reset() is equal to the number of items in the queue.
//
// The circular buffer needs all elements in the queue to exist to simplify
// the enqueue operation..
func (c *Circular) zeroQueue() int {
	var x int
	for i := len(c.Items); i < cap(c.Items); i++ {
		c.Items = append(c.Items, nil)
		x++
	}
	return x
}
