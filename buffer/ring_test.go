package buffer

import (
  "testing"
)

func TestRingBuffer(t *testing.T) {
  r := NewRing(5)
  if r.Size() != 5 {
        t.Errorf("Expected ring buffers cap to be 5, got %d", r.Size())
  }
}

func TestRingBufferEnqueueDequeuePeek(t * testing.T) {
  tests := []struct{
    size int
    items []int
    tailPos int
    dequeue []int
    dequeueRemaining []int
    dequeueHeadPos int
    dequeueTailPos int
    dequeueIsEmpty bool
    enqueue1 []int
    enqueue1Remaining []int
    enqueue1HeadPos int
    enqueue1TailPos int
    enqueueIsEmpty bool
    enqueue2 []int
    enqueue2Remaining []int
    enqueue2HeadPos int
    enqueue2TailPos int
  }{
    {4, []int{0, 1, 2}, 3, []int{0, 1}, []int{2}, 2, 3, false, []int{3, 4},
      []int{2, 3, 4}, 2, 0, false, []int{}, []int{2,3,4}, 2, 0},
    {4, []int{0, 1, 2, 3}, 4, []int{0, 1, 2}, []int{3}, 3, 4, false, []int{4, 5},
      []int{3, 4, 5}, 3, 1, false, []int{6, 7, 8}, []int{5, 6, 7, 8}, 0, 4},
  }
  for i, test := range tests {
    b := NewRing(test.size)
    for _, v := range test.items {
      _ = b.Enqueue(v)
    }
    if b.Tail != test.tailPos {
      t.Errorf("%d: post initial enqueue, expected tail to be at pos %d, was at %d", i, test.tailPos, b.Tail)
    }
    for j, v := range test.dequeue {
      val, _ := b.Dequeue()
      if v != val {
        t.Errorf("%d dequeue item %d: expected %v got %v", i, j, v, val)
      }
    }
    if b.Head != test.dequeueHeadPos {
      t.Errorf("%d: post dequeue, expected head pos to be %d, got %d", i, test.dequeueHeadPos, b.Head)
    }
    if b.Tail != test.dequeueTailPos {
      t.Errorf("%d: post dequeue, expected Tail pos to be %d, got %d", i, test.dequeueTailPos, b.Tail)
    }
    if b.Len() != len(test.dequeueRemaining) {
      t.Errorf("%d: after dequeue, expected %d items in buffer, got %d", i, len(test.dequeueRemaining), b.Len())
      goto Enqueue1
    }
    for j, v := range test.dequeueRemaining {
      k := b.Head + j
      if k >= cap(b.Items) {
        k -= cap(b.Items)
      }
      if v != b.Items[k] {
        t.Errorf("%d dequeueRemainingItem %d: expected %v, got %v", i, j, v, b.Items[k])
      }
    }
Enqueue1:
    for _, v := range test.enqueue1 {
      _ = b.Enqueue(v)
    }
    if b.Head != test.enqueue1HeadPos {
      t.Errorf("%d: after enqueue1, expected head pos to be %d got %d", i, test.enqueue1HeadPos, b.Head)
    }
    if b.Tail != test.enqueue1TailPos {
      t.Errorf("%d: after enqueue1, expected head pos to be %d got %d", i, test.enqueue1TailPos, b.Tail)
    }
    if b.Len() != len(test.enqueue1Remaining) {
      t.Errorf("%d: after enqueue1, expected %d items in buffer, got %d", i, len(test.enqueue1Remaining), b.Len())
      goto Enqueue2
    }
    for j, v := range test.enqueue1Remaining {
      k := b.Head + j
      if k >= cap(b.Items) {
        k -= cap(b.Items)
      }
      if v != b.Items[k] {
        t.Errorf("%d enqueue1RemainingItem %d: expected %v, got %v", i, j, v, b.Items[k])
      }
    }
Enqueue2:
    for _, v := range test.enqueue2 {
      _ = b.Enqueue(v)
    }
    if b.Head != test.enqueue2HeadPos {
      t.Errorf("%d: after enqueue2, expected head pos to be %d got %d", i, test.enqueue2HeadPos, b.Head)
    }
    if b.Tail != test.enqueue2TailPos {
      t.Errorf("%d: after enqueue2, expected head pos to be %d got %d", i, test.enqueue2TailPos, b.Tail)
    }
    if b.Len() != len(test.enqueue2Remaining) {
      t.Errorf("%d: after enqueue2, expected %d items in buffer, got %d", i, len(test.enqueue2Remaining), b.Len() )
      goto Enqueue2
    }
    for j, v := range test.enqueue2Remaining {
      k := b.Head + j
      if k >= cap(b.Items) {
        k -= cap(b.Items)
      }
      if v != b.Items[k] {
        t.Errorf("%d enqueue2RemainingItem %d: expected %v, got %v", i, j, v, b.Items[k])
      }
    }
  }
}

func TestBufferPos(t *testing.T) {
  tests := []struct {
      val string
      dequeue bool
      ok bool
      expectedHead int
      expectedTail int
  }{
    {"", true, false, 0, 0},
    {"a", false, false, 0, 1},
    {"b", false, false, 0, 2},
    {"c", false, false, 0, 3},
    {"d", false, false, 0, 4},
    {"e", false, false, 0, 5},
    {"f", false, false, 0, 6},
    {"a", true, true, 1, 6},
    {"b", true, true, 2, 6},
    {"g", false, false, 2, 0},
    {"h", false, false, 2, 1},
    {"i", false, false, 3, 2},
    {"j", false, false, 4, 3},
    {"k", false, false, 5, 4},
    {"l", false, false, 6, 5},
    {"g", true, true, 0, 5},
    {"h", true, true, 1, 5},
    {"i", true, true, 2, 5},
    {"j", true, true, 3, 5},
    {"k", true, true, 4, 5},
    {"l", true, true, 5, 5},
    {"", true, false, 5, 5},
  }
  b := NewRing(6)
  for i, test := range tests {
    if test.dequeue{
      v, ok := b.Dequeue()
      if test.ok != ok {
        t.Errorf("%d: dequeue expected %t, got %t", test.ok, ok)
      }
      if test.ok {
        if v != test.val {
          t.Errorf("%d: dequeue val expected to be %s, got %s", i, test.val, v)
        }
        if b.Head != test.expectedHead {
          t.Errorf("%d: post dequeue expected head to be at pos %d, was at %d", i, test.expectedHead, b.Head)
        }
        if b.Tail != test.expectedTail {
          t.Errorf("%d: post dequeue expected tail to be at pos %d, was at %d", i, test.expectedTail, b.Tail)
        }
      }
      continue
    }
    _ = b.Enqueue(test.val)
    if b.Head != test.expectedHead {
      t.Errorf("%d: post enqueue, expected head to be at pos %d, was at %d", i, test.expectedHead, b.Head)
    }
    if b.Tail != test.expectedTail {
      t.Errorf("%d: post enqueue, expected head to be at pos %d, was at %d", i, test.expectedTail, b.Tail)
    }
  }
}
