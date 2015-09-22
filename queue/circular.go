package dq

import (
  "fmt"
  "math"
  "sync"
)
// Circular is a bounded queue implemented as a circular queue.
type Circular struct {
  sync.Mutex
  items []interface{}
  head int
  tail int
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
  return &Circular{items: make([]interface{}, size + 1, size + 1), cap: size}
}

// Enqueue will return an error if the queue is full
func (c *Circular) Enqueue(item interface{}) error {
  c.Lock()
  if c.isFull() {
    c.Unlock()
    return fmt.Errorf("queue full: cannot enqueue %v", item)
  }
  c.items[c.tail] = item
  if c.tail == c.cap {
    c.tail = 0
  } else {
    c.tail++
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
  item := c.items[c.head]
  if c.head == c.cap {
    c.head = 0
  } else {
    c.head++
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
  return c.items[c.head], true
}

// IsEmpty returns whether or not the queue is empty
func (c *Circular) IsEmpty() bool {
  c.Lock()
  if c.head == c.tail {
    c.Unlock()
    return true
  }
  c.Unlock()
  return false
}

// isEmpty is an unexported version that expects the caller to handle locking.
// This eliminates double locking on dequeue and peek
func (c *Circular) isEmpty() bool {
  if c.head == c.tail {
    return true
  }
  return false
}

// IsFull returns whether or not the queue is full
func (c *Circular) IsFull() bool {
  c.Lock()
  if c.head != int(math.Mod(float64(c.tail + 1), float64(cap(c.items)))) {
    c.Unlock()
    return false
  }
  c.Unlock()
  return true
}

// isFull is an unexported version that expects the caller to handle locking.
// This eliminates double locking on enqueue
func (c *Circular) isFull() bool {
  if c.head != int(math.Mod(float64(c.tail + 1), float64(cap(c.items)))) {
    return false
  }
  return true
}
