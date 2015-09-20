package dq
import (
	"testing"
)

func TestNew(t *testing.T) {
	q := NewQ(10, false)
	if cap(q.items) != 10 {
		t.Errorf("expected 10, got %d", cap(q.items))
	}
	if q.bounded {
		t.Errorf("expected false, got %t", q.bounded)
	}

	q = NewQueue(100, true)
	if cap(q.items) != 100 {
		t.Errorf("expected 100, got %d", cap(q.items))
	}
	if !q.bounded {
		t.Errorf("expected true, got %t", q.bounded)
	}
}

// tests enqueue, growth, capacity restriction, and basic dequeue
func TestQueueing(t *testing.T) {
	var tests = []struct {
		size        int
		bounded     bool
		headPos     int
		expectedLen     int
		expectedCap int
		items       []interface{}
		errString   string
	}{
		{size: 2, bounded: false, expectedLen: 4, expectedCap: 4, items: []interface{}{0, 1, 2, 3}, errString: ""},
		{size: 2, bounded: false, expectedLen: 5, expectedCap: 8, items: []interface{}{0, 1, 2, 3, 4}, errString: ""},
		{size: 4, bounded: true, expectedLen: 4, expectedCap: 4, items: []interface{}{0, 1, 2, 3}, errString: ""},
	}
	for i, test := range tests {
		q := NewQ(test.size, test.bounded)
		for _, v := range test.items {
			_ = q.Enqueue(v)

		}

		// check that the items are as expected:
		if len(q.items) != test.expectedLen {
			t.Errorf("%d: expected %d items in queue, got %d", i, test.expectedLen, len(q.items))
		}
		if cap(q.items) != test.expectedCap {
			t.Errorf("%d: expected queue cap to be %d, got %d", i, test.expectedCap, cap(q.items))
		}
		if q.head != test.headPos {
			t.Errorf("%d: expected head to be at pos %d, got %d", i, test.headPos, q.head)
		}
		if q.bounded != test.bounded {
			t.Errorf("%d: expected bounded to be %t, was %t", i, test.bounded, q.bounded)
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
func TestDequeueEnqueue(t *testing.T) {
	tests := []struct {
		size        int
		bounded     bool
		headPos     int
		expectedLen int
		expectedCap int
		dequeueCnt  int
		dequeueVals []interface{}
		items       []interface{}
		items2      []interface{}
		errString   string
	}{
		{size: 10, bounded: false, expectedLen: 7, expectedCap: 10, dequeueCnt: 5, dequeueVals: []interface{}{0, 1, 2, 3, 4},
			items: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, items2: []interface{}{10, 11}, errString: ""},
	}

	// First add the queue
	for _, test := range tests {
		var err error
		q := NewQ(test.size, test.bounded)
		for _, v := range test.items {
			err = q.Enqueue(v)
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
		if len(q.items) != test.expectedLen {
			t.Errorf("Expected tail to be at %d, got %d", test.expectedLen, len(q.items))
		}
		if cap(q.items) != test.expectedCap {
			t.Errorf("Expected cap of queue to be %d. got %d", test.expectedCap, cap(q.items))
		}

	}
}

func TestSetShiftPercentage(t *testing.T) {
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
	q := NewQueue(10, true)
	for i, test := range tests {
		q.SetShiftPercent(test.percent)
		if q.shiftPercent != test.expected {
			t.Errorf("%d: expected shiftPercent to be %d; got %d", i, test.expected, q.shiftPercent)
		}
	}
}

func TestCappedQueue(t *testing.T) {
	q := NewQ(4, true)
	for i := 0; i < 4; i++ {
		q.Enqueue(i)
	}
	// remove an item and then try to add
	_ = q.Dequeue()
	err := q.Enqueue(5)
	if err != nil {
		t.Errorf("Expected enqueue to a bounded queue with len == cap, but room to shift to succeed; got %s", err)
		return
	}
	// enqueue another item, queue is full, this should fail
	err = q.Enqueue(6)
	if err == nil {
		t.Errorf("Expected enqueue to a bounded queue that is full to error, it did not")
		return
	}
	if err.Error() != "bounded queue full: cannot enqueue '6'" {
		t.Errorf("Expected enqueue to a capped queue to error with \"bounded queue full: cannot enqueue '6'\" , got %q", err)
	}

}
