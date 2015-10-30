package queue

// Copyright 2015 by Joel Scoble
//
// This implementation of a priority queue using a heap is based on the code
// provided in https://http://golang.org/src/container/heap/example_pq_test.go
//
// This implemantion uses interface{} and is thread safe.
//
// The original copyright notice:
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"container/heap"
	"sync"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; arbitrary.
	priority int         // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A HeapPriority implements heap.Interface and holds Items.
type HeapPriority struct {
	mu    *sync.Mutex
	items PQueue
}

// PQueue represents a priority queue
type PQueue []*Item

func (pq PQueue) Len() int { return len(pq) }

func (pq PQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push pushes an item onto the priority queue.
func (pq *PQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop pops the next itme from the priority queue.
func (pq *PQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and vlaue of an Item in the queue.
func (pq *PQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// NewHeapPriority returns a new priority queue with the item's cap set at l; if l > 0.
func NewHeapPriority(l int) *HeapPriority {
	if l <= 0 {
		return &HeapPriority{}
	}
	return &HeapPriority{items: make([]*Item, l, l)}
}

func (pq HeapPriority) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.items.Len()
}

func (pq HeapPriority) Less(i, j int) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.items.Less(i, j)
}

func (pq HeapPriority) Swap(i, j int) {
	pq.mu.Lock()
	pq.items.Swap(i, j)
	pq.mu.Unlock()
}

// Push pushes an item onto the priority queue.
func (pq *HeapPriority) Push(x interface{}) {
	pq.mu.Lock()
	pq.items.Push(x)
	pq.mu.Unlock()
}

// Pop pops the next item from the priority queue.
func (pq *HeapPriority) Pop() interface{} {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.items.Pop()
}

// update modifies the priority and value of an Item in the queue.
func (pq *HeapPriority) update(item *Item, value string, priority int) {
	pq.mu.Lock()
	pq.items.update(item, value, priority)
	pq.mu.Unlock()
}
