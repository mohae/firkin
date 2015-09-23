package queue

import (
  "testing"
)

func TestCircularQ(t *testing.T) {
  tests := []struct{
    size int
    items []int
    expectedQueue int
    initHead int
    initTail int
    initIsFull bool
    initIsEmpty bool
    dequeue []int
    dequeueOk bool
    dequeueHead int
    dequeueTail int
    dequeueIsFull bool
    dequeueIsEmpty bool
    dequeueLen int
    enqueue []int
    enqueueHead int
    enqueueTail int
    enqueueIsFull bool
    enqueueIsEmpty bool
    enqueueLen int
    err string
  }{
    {2, []int{}, 3, 0, 0, false, true, []int{}, false, 0, 0, false, true, 0, []int{}, 0, 0, false, true, 0, ""},
    {2, []int{0}, 3, 0, 1, false, false, []int{0}, true, 1, 1, false, true, 0, []int{1, 2}, 1, 0, true, false, 2, ""},
    {2, []int{0, 1}, 3, 0, 2, true, false, []int{0, 1}, true, 2, 2, false, true, 0, []int{2, 3}, 2, 1, true, false, 2, ""},
    {2, []int{0, 1}, 3, 0, 2, true, false, []int{0, 1}, true, 2, 2, false, true, 0, []int{2, 3, 4}, 2, 1, true, false, 2, "queue full: cannot enqueue 4"},
    {4, []int{0, 1, 2, 3}, 5, 0, 4, true, false, []int{0, 1, 2}, true, 3, 4, false, false, 1, []int{4, 5, 6}, 3, 2, true, false, 4, ""},
    {4, []int{0, 1, 2, 3}, 5, 0, 4, true, false, []int{0, 1, 2, 3}, true, 4, 4, false, true, 0, []int{}, 4, 4, false, true, 0, ""},
  }
  for i, test := range tests {
    cq := NewCircularQ(test.size)
    for _, v := range test.items {
      _ = cq.Enqueue(v)
    }
    if cq.Head != test.initHead {
      t.Errorf("%d initial: expected Head to be %d, got %d", i, test.initHead, cq.Head)
    }
    if cq.Tail != test.initTail {
      t.Errorf("%d initial: expected Tail to be %d, got %d", i, test.initTail, cq.Tail)
    }
    if cq.IsEmpty() != test.initIsEmpty {
      t.Errorf("%d initial: expected isEmpty to be %t, got %t", i, test.initIsEmpty, cq.IsEmpty())
    }
    if cq.IsFull() != test.initIsFull {
      t.Errorf("%d initial: expected isFull to be %t, got %t", i, test.initIsFull, cq.IsFull())
    }
    for j, v := range test.dequeue{
      val, ok := cq.Dequeue()
      if ok != test.dequeueOk {
        t.Errorf("%d dequeue #%d: expected %t, got %t", i, j, test.dequeueOk, ok)
      }
      if val != v {
        t.Errorf("%d: dequeue item %d: expected %v got %v", i, j, v, val)
      }
    }
    if cq.Head != test.dequeueHead {
      t.Errorf("%d dequeue: expected Head to be %d, got %d", i, test.dequeueHead, cq.Head)
    }
    if cq.Tail != test.dequeueTail {
      t.Errorf("%d dequeue: expected Tail to be %d, got %d", i, test.dequeueTail, cq.Tail)
    }
    if cq.IsEmpty() != test.dequeueIsEmpty {
      t.Errorf("%d dequeue: expected isEmpty to be %t, got %t", i, test.dequeueIsEmpty, cq.IsEmpty())
    }
    if cq.IsFull() != test.dequeueIsFull {
      t.Errorf("%d dequeue: expected isFull to be %t, got %t", i, test.dequeueIsFull, cq.IsFull())
    }
    if cq.Len() != test.dequeueLen {
      t.Errorf("%d dequeue: expected len to be %d, got %d", i, test.dequeueLen, cq.Len())
    }
    var err error
    for j, v := range test.enqueue {
      err = cq.Enqueue(v)
      if err != nil {
        if err.Error() != test.err {
          t.Errorf("%d enqueue #%d: expected error to be %q, got %q", i, j, test.err, err.Error())
        }
      }
    }
    if err == nil && test.err != "" {
      t.Errorf("%d enqueue: expected error an error: %q, got none", i, test.err)
    }
    if cq.Head != test.enqueueHead {
      t.Errorf("%d enqueue: expected Head to be %d, got %d", i, test.enqueueHead, cq.Head)
    }
    if cq.Tail != test.enqueueTail {
      t.Errorf("%d enqueue: expected Tail to be %d, got %d", i, test.enqueueTail, cq.Tail)
    }
    if cq.IsEmpty() != test.enqueueIsEmpty {
      t.Errorf("%d enqueue: expected isEmpty to be %t, got %t", i, test.enqueueIsEmpty, cq.IsEmpty())
    }
    if cq.IsFull() != test.enqueueIsFull {
      t.Errorf("%d enqueue: expected isFull to be %t, got %t", i, test.enqueueIsFull, cq.IsFull())
    }
    if cq.Len() != test.enqueueLen {
      t.Errorf("%d enqueue: expected len to be %d, got %d", i, test.enqueueLen, cq.Len())
    }
  }
}
