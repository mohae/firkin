package queue

import (
  "fmt"
  "math"
  "sync"
)
// Circular is a bounded queue implemented as a circular queue.  Even though
// Items, Head, and Tail are exported, in most cases, they should not be
// directly.  Doing so may lead to outcomes less than desirable. Use the
// exported methods to interact with the Circular queue.
type Circular struct {
  sync.Mutex
  Items []interface{}
  Head int
  Tail int
  cap int
}

// NewCircularQ returns an initialized circular queue. Even though creating
// the slice with an initial length is much slower than creating one without
// the initial length, cap only, this is done to simplify the actual queue
// management. Don't need to worry about appending vs adding via index and
// don't need to check to see if an append will cause the slice to grow.
//
// The slice is 1 slot larger than the requested size for empty/full
// detection.
func NewCircularQ(size int) *Circular {
  return &Circular{Items: make([]interface{}, size + 1, size + 1), cap: size}
}

// Enqueue will return an error if the queue is full
func (c *Circular) Enqueue(item interface{}) error {
  c.Lock()
  if c.isFull() {
    c.Unlock()
    return fmt.Errorf("queue full: cannot enqueue %v", item)
  }
  c.Items[c.Tail] = item
  if c.Tail == c.cap {
    c.Tail = 0
  } else {
    c.Tail++
  }
  c.Unlock()
  return nil
}

// Dequeue will remove an item from the queue and return it. If the queue is
// empty, a false will be returned.
func (c *Circular) Dequeue() (interface{}, bool) {
  c.Lock()
  if c.isEmpty() {
    c.Unlock()
    return nil, false
  }
  item := c.Items[c.Head]
  if c.Head == c.cap {
    c.Head = 0
  } else {
    c.Head++
  }
  c.Unlock()
  return item, true
}

// Peek will return the next item in the queue without removing it from the
// queue. If the queue is empty, a false will be returned.
func (c *Circular) Peek() (interface{}, bool) {
  c.Lock()
  defer c.Unlock()
  if c.isEmpty() {
    return nil, false
  }
  return c.Items[c.Head], true
}

// IsEmpty returns whether or not the queue is empty
func (c *Circular) IsEmpty() bool {
  c.Lock()
  if c.Head == c.Tail {
    c.Unlock()
    return true
  }
  c.Unlock()
  return false
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
  if c.Head != int(math.Mod(float64(c.Tail + 1), float64(cap(c.Items)))) {
    c.Unlock()
    return false
  }
  c.Unlock()
  return true
}

// isFull is an unexported version that expects the caller to handle locking.
// This eliminates double locking on enqueue
func (c *Circular) isFull() bool {
  if c.Head != int(math.Mod(float64(c.Tail + 1), float64(cap(c.Items)))) {
    return false
  }
  return true
}

// returns the Size of the Queue
func (c *Circular) Size() int {
  return c.cap
}
