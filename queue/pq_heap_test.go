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
package queue

import (
	"container/heap"
	"testing"
)

func TestPQHeap(t *testing.T) {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := NewHeapPriority(len(items))
	i := 0
	for value, priority := range items {
		pq.items[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq.items)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	pq.Push(item)
	pq.update(item, "grapefruit", 5)

	// Take the items out; they arrive in decreasing priority order.
	expected := []struct {
		priority int
		value    string
	}{
		{5, "grapefruit"},
		{4, "pear"},
		{3, "banana"},
		{2, "apple"},
	}
	i = 0
	for pq.Len() > 0 {
		item := heap.Pop(&pq.items).(*Item)
		if item.priority != expected[i].priority || item.value != expected[i].value {
			t.Errorf("Expected %v got %v", expected[i], item)
		}
		i++
	}
}
