package stack

import (
	"fmt"
	"sync"
)

type Stack struct {
	rw      sync.RWMutex
	items   []interface{}
	cap     int
	size    int
	bounded bool
}

// NewStack returns a new stack with its initial capacity equal to the received
// size and bounded set accordingly.
func NewStack(cap int, bounded bool) *Stack {
	return &Stack{items: make([]interface{}, 0, cap), cap: cap, bounded: bounded}
}

// Push an item on the stack. An error will occur is the stack is bounded
// and at capacity.
func (s *Stack) Push(item interface{}) error {
	s.rw.Lock()
	if s.bounded && s.size == s.cap {
		s.rw.Unlock()
		return fmt.Errorf("bounded stack full: cannot push '%v' onto the stack", item)
	}
	if s.size == len(s.items) {
		s.items = append(s.items, item)
		s.size++
		s.rw.Unlock()
		return nil
	}
	s.items[s.size] = item
	s.size++
	s.rw.Unlock()
	return nil
}

// Pop pops an item off the stack {}. A nil wil be returned if the stack is
// empty
func (s *Stack) Pop() (interface{}, bool) {
	s.rw.Lock()
	defer s.rw.Unlock()
	if s.size == 0 {
		return nil, false
	}
	s.size--
	return s.items[s.size], true
}

// Peek returns the item at the top of the stack without popping it. If the
// stack is empty, it will return nil
func (s *Stack) Peek() (interface{}, bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if s.size == 0 {
		return nil, false
	}
	return s.items[s.size-1], true
}

// IsEmpty returns whether or not the stack is empty
func (s *Stack) IsEmpty() bool {
	s.rw.RLock()
	if s.size == 0 {
		s.rw.RUnlock()
		return true
	}
	s.rw.RUnlock()
	return false
}

// Size returns the current size of the stack (number of items)
func (s *Stack) Size() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.size
}

// Reset restets the stack: the capacity of the stack will be reset to its
// initial capacity. Anything in the queue will be lost
func (s *Stack) Reset() {
	s.rw.Lock()
	s.size = 0
	s.items = make([]interface{}, 0, s.cap)
	s.rw.Unlock()
}
