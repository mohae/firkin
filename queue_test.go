package dq
import (
	"testing"
)

func TestNew(t *testing.T) {
	q := NewQ(10)
	if q.Cap() != 10 {
		t.Errorf("expected 10, got %d", cap(q.items))
	}
	q = NewQueue(100)
	if q.Cap() != 100 {
		t.Errorf("expected 100, got %d", cap(q.items))
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
			t.Errorf("%d: expected %d items in queue, got %d", i, test.expectedLen, len(q.items))
		}
		if q.Cap() != test.expectedCap {
			t.Errorf("%d: expected queue cap to be %d, got %d", i, test.expectedCap, cap(q.items))
		}
		if q.head != test.headPos {
			t.Errorf("%d: expected head to be at pos %d, got %d", i, test.headPos, q.head)
		}
		for j := 0; j < len(q.items); j++ {
			if q.items[j] != test.items[j] {
				t.Errorf("%d: expected value of index %d to be %d, got %d", i, j, test.items[j], q.items[j])
			}
		}

		// dequeue 1 item and check
		next := q.Dequeue()
		if next != test.items[0] {
			t.Errorf("%d: expected %d, got %d", i, test.items[0], next)
			continue
		}
		if q.head != 1 {
			t.Errorf("%d: expected head to point to 1, got %d", i, q.head)
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
		dequeueCnt  int
		dequeueVals []interface{}
		items       []interface{}
		items2      []interface{}
		errString   string
	}{
		{size: 10, expectedLen: 7, expectedCap: 10, dequeueCnt: 5, dequeueVals: []interface{}{0, 1, 2, 3, 4},
			items: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, items2: []interface{}{10, 11}, errString: ""},
	}

	// First add the queue
	for _, test := range tests {
		var err error
		q := NewQueue(test.size)
		for _, v := range test.items {
			_ = q.Enqueue(v)
		}
		if test.errString != "" {
			if err == nil {
				t.Errorf("Expected error, got none")
				continue
			}
			if err.Error() != test.errString {
				t.Errorf("Expected error to be %q. got %q", test.errString, err.Error())
				continue
			}
		}

		// dequeue 5 items
		for i := 0; i < test.dequeueCnt; i++ {
			v := q.Dequeue()
			if v != test.dequeueVals[i] {
				t.Errorf("Expected %v, got %v", test.dequeueVals[i], v)
			}
		}

		if q.head != test.dequeueCnt {
			t.Errorf("Expected head to point to %d, got %d", test.dequeueCnt, q.head)
		}
		// enqueue the next items; should not grow, should just shift the items
		for _, v := range test.items2 {
			q.Enqueue(v)
		}
		if q.head != 0 {
			t.Errorf("Expected head to be at pos 0, got %d", q.head)
		}
		if q.Len() != test.expectedLen {
			t.Errorf("Expected tail to be at %d, got %d", test.expectedLen, len(q.items))
		}
		if q.Cap() != test.expectedCap {
			t.Errorf("Expected cap of queue to be %d. got %d", test.expectedCap, cap(q.items))
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

func TestQIsEmpty(t *testing.T) {
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
	}
}
