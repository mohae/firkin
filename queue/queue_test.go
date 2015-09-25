package queue
import (
	"testing"
)

func TestNew(t *testing.T) {
	q := NewQ(10)
	if q.Cap() != 10 {
		t.Errorf("expected 10, got %d", cap(q.Items))
	}
	q = NewQueue(100)
	if q.Cap() != 100 {
		t.Errorf("expected 100, got %d", cap(q.Items))
	}
}

// tests enqueue, growth, capacity restriction, and basic dequeue
func TestQueueing(t *testing.T) {
	var tests = []struct {
		size        int
		headPos     int
		expectedLen int
		expectedCap int
		items       []interface{}
	}{
		{size: 2, expectedLen: 2, expectedCap: 2, items: []interface{}{0, 1}},
		{size: 2, expectedLen: 5, expectedCap: 8, items: []interface{}{0, 1, 2, 3, 4}},
		{size: 2, expectedLen: 4, expectedCap: 4, items: []interface{}{0, 1, 2, 3}},
	}
	for i, test := range tests {
		q := NewQ(test.size)
		for _, v := range test.items {
			_ = q.Enqueue(v)

		}
		// check that the items are as expected:
		if q.Len() != test.expectedLen {
			t.Errorf("%d: expected %d items in queue, got %d", i, test.expectedLen, len(q.Items))
		}
		if q.Cap() != test.expectedCap {
			t.Errorf("%d: expected queue cap to be %d, got %d", i, test.expectedCap, cap(q.Items))
		}
		if q.Head != test.headPos {
			t.Errorf("%d: expected head to be at pos %d, got %d", i, test.headPos, q.Head)
		}
		for j := 0; j < len(q.Items); j++ {
			if q.Items[j] != test.items[j] {
				t.Errorf("%d: expected value of index %d to be %d, got %d", i, j, test.items[j], q.Items[j])
			}
		}

		// dequeue 1 item and check
		next, _ := q.Dequeue()
		if next != test.items[0] {
			t.Errorf("%d: expected %d, got %d", i, test.items[0], next)
			continue
		}
		if q.Head != 1 {
			t.Errorf("%d: expected head to point to 1, got %d", i, q.Head)
		}
	}
}

// Tests Enqueue/Dequeue/Enqueue, shifting, and growth is properly handled
func TestQDequeueEnqueue(t *testing.T) {
	tests := []struct {
		size        int
		headPos     int
		expectedLen int
		expectedCap int
		expectedPeek int
		postPeekHeadPos int
		dequeueCnt  int
		dequeueVals []interface{}
		postDequeueHeadPos int
		postDequeueLen int
		items       []interface{}
		enqueueItems      []interface{}
		postEnqueueLen int
		postEnqueueCap int
	}{
		{size: 10, headPos: 0, expectedLen: 10, expectedCap: 10, expectedPeek: 5, postPeekHeadPos: 5,
		 	dequeueCnt: 5, dequeueVals: []interface{}{0, 1, 2, 3, 4}, postDequeueHeadPos: 5,
			postDequeueLen: 5, items: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			enqueueItems: []interface{}{10, 11}, postEnqueueLen: 7, postEnqueueCap: 10},
	}

	// First add the queue
	for i, test := range tests {
		q := NewQueue(test.size)
		for _, v := range test.items {
			_ = q.Enqueue(v)
		}
		if q.Head != test.headPos {
			t.Errorf("%d: post queue population, expected head pos to be %d, got %d", test.headPos, q.Head)
		}
		if q.Len() != test.expectedLen {
			t.Errorf("%d: post queue population, expected len to be %d got %d", i, test.expectedLen, q.Len())
		}
		if q.Cap() != test.expectedCap {
			t.Errorf("%d: post queue population, expected cap to be %d, got %d", i, test.expectedCap, q.Cap())
		}
		// dequeue 5 items
		for i := 0; i < test.dequeueCnt; i++ {
			v , _ := q.Dequeue()
			if v != test.dequeueVals[i] {
				t.Errorf("%d: dequeue: expected %v, got %v", i,test.dequeueVals[i], v)
			}
		}
		if q.Head != test.dequeueCnt {
			t.Errorf("%d: post deuque: expected head to point to %d, got %d", i, test.dequeueCnt, q.Head)
		}
		// peek stuff
		v, _ := q.Peek()
		if v.(int) != test.expectedPeek {
			t.Errorf("%d: post peek: expected peek to return %d, got %d", i, test.expectedPeek, v.(int))
		}
		if q.Head != test.postPeekHeadPos {
			t.Errorf("%d: post peek: expected head to be at pos %d, got %d", i, test.postPeekHeadPos, q.Head)
		}
		// enqueue the next items; should not grow, should just shift the items
		for _, v := range test.enqueueItems {
			q.Enqueue(v)
		}
		if q.Head != 0 {
			t.Errorf("%d post enqueue: expected head to be at pos 0, got %d", i, q.Head)
		}
		if q.Len() != test.postEnqueueLen {
			t.Errorf("%d post enqueue: expected tail to be at %d, got %d", i, test.postEnqueueLen, q.Len())
		}
		if q.Cap() != test.postEnqueueCap {
			t.Errorf("%d post enqueue: expected cap of queue to be %d, got %d", i, test.postEnqueueCap, q.Cap())
		}
	}
}

func TestQSetShiftPercentage(t *testing.T) {
	tests := []struct{
		percent int
		expected int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{20, 20},
		{99, 99},
		{100, 100},
		{101, 100},
	}
	q := NewQueue(10)
	for i, test := range tests {
		q.SetShiftPercent(test.percent)
		if q.shiftPercent != test.expected {
			t.Errorf("%d: expected shiftPercent to be %d; got %d", i, test.expected, q.shiftPercent)
		}
	}
}

func TestQIsEmptyFull(t *testing.T) {
	tests := []struct{
		size int
		items []int
		isEmpty bool
	}{
		{4, []int{}, true},
		{4, []int{0, 1, 2, 3}, false},
	}
	for i, test := range tests {
		q := NewQ(test.size)
		for _, v := range test.items {
			q.Enqueue(v)
		}
		if q.IsEmpty() != test.isEmpty {
			t.Errorf("%d: expected IsEmpty() to return %t. got %t", i, test.isEmpty, q.IsEmpty())
		}
		if q.IsFull() {
			t.Errorf("%d: expected IsFull() to return false, got %t", i, q.IsFull())
		}
	}
}

func TestDequeuePeekErr(t *testing.T) {
	tests := []struct{
		size int
		items []int
		retItems []int
		retOk []bool
		isEmpty bool
		isFull bool
	}{
		{2, []int{0, 1, 2, 3, 4}, []int{}, []bool{}, false, false},
		{2, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}, []bool{true, true, true, true, true}, true, false},
		{2, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4, 5}, []bool{true, true, true, true, true, false}, true, false},
	}
	for i, test := range tests {
		q := NewQ(test.size)
		for j, v := range test.items {
			err := q.Enqueue(v)
			if err != nil {
				t.Errorf("%d enqueueing #%d: unexpected error %q", i, j)
			}
		}
		for j, v := range test.retItems {
			val, ok := q.Peek()
			if ok != test.retOk[j] {
				t.Errorf("%d peek #%d: expected peek to return %t, got %t", i, j, test.retOk[j], ok)
			}
			if ok {
				if val != v {
					t.Errorf("%d peek #%d: expected peek to return %v, got %v", i, j, v, val)
				}
			}
			val, ok = q.Dequeue()
			if ok != test.retOk[j] {
				t.Errorf("%d peek #%d: expected peek to return %t, got %t", i, j, test.retOk[j], ok)
			}
			if ok {
				if val != v {
					t.Errorf("%d peek #%d: expected peek to return %v, got %v", i, j, v, val)
				}
			}
		}
		if q.IsEmpty() != test.isEmpty {
			t.Errorf("%d: expected queue IsEmpty to be %t, got %t", i, test.isEmpty, q.IsEmpty())
		}
		if q.IsFull() != test.isFull {
			t.Errorf("%d: expected queue IsFull to be %t, got %t", i, test.isFull, q.IsFull())
		}
	}
}

func TestQueueResetResize(t *testing.T) {
	tests := []struct{
		size int
		enqueue int
		dequeue int
		cap int
		resize int
		expectedLen int
		expectedCap int
	}{
	  {4, 0, 0, 4, 0, 0, 4},
		{2, 2, 0, 2, 0, 2, 2},
		{2, 2, 2, 2, 0, 0, 2},
		{4, 2, 2, 4, 0, 0, 4},
		{4, 2, 0, 4, 0, 2, 4},
		{4, 1, 0, 4, 0, 1, 4},
		{4, 2, 1, 4, 0, 1, 4},
		{2, 5, 0, 8, 0, 5, 6},
		{2, 5, 5, 8, 0, 0, 2},
		{2, 5, 5, 8, 4, 0, 4},
		{2, 6, 1, 8, 0, 5, 6},
		{2, 6, 1, 8, 3, 5, 6},
		{2, 6, 1, 8, 7, 5, 7},
	}
	for i, test := range tests {
		q := NewQ(test.size)
		for j := 0; j < test.enqueue; j++ {
			_ = q.Enqueue(j)
		}
		if q.Len() != test.enqueue {
			t.Errorf("%d: expected queue len to be %d, got %d", i, test.enqueue, q.Len())
		}
		if q.Cap() != test.cap {
			t.Errorf("%d: expected queue cap to be %d, got %d", i, test.cap, q.Cap())
		}
		for j := 0; j < test.dequeue; j++ {
			_, _ = q.Dequeue()
		}
		q.Reset()
		if q.Len() != 0 {
			t.Errorf("%d: after Reset(), expected queue len to be 0, got %d", i, q.Len())
		}
		if q.Head != 0 {
			t.Errorf("%d: after Reset(), expected queue head to be at pos 0, was at pos %d", i, q.Head)
		}
		if q.Cap() != test.cap {
			t.Errorf("%d: after Reset(), expected queue cap to be %d, got %d", i, test.cap, q.Cap())
		}
	}

	for i, test := range tests {
		q := NewQ(test.size)
		for j := 0; j < test.enqueue; j++ {
			_ = q.Enqueue(j)
		}
		if q.Len() != test.enqueue {
			t.Errorf("%d: expected queue len to be %d, got %d", i, test.enqueue, q.Len())
		}
		if q.Cap() != test.cap {
			t.Errorf("%d: expected queue cap to be %d, got %d", i, test.cap, q.Cap())
		}
		for j := 0; j < test.dequeue; j++ {
			_, _ = q.Dequeue()
		}
		q.Resize(test.resize)
		if q.Len() != test.expectedLen{
			t.Errorf("%d: after Resize(), expected queue len to be %d, got %d", i, test.expectedLen,  q.Len())
		}
		if q.Head != 0 {
			t.Errorf("%d: after Resize(), expected queue head to be at pos 0, was at pos %d", i, q.Head)
		}
		if q.Cap() != test.expectedCap {
			t.Errorf("%d: after Resize(), expected queue cap to be %d, got %d", i, test.expectedCap, q.Cap())
		}
	}
}
