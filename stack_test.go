package dq

import (
  "testing"
)

func TestPushStack(t *testing.T) {
  tests := []struct{
    cap int
    bounded bool
    initVals []int
    initSize int
    pushVals []int
    pushedSize int
    pop bool
    popVal interface{}
    poppedSize int
    peekVal interface{}
    expectedErr string
    isEmpty bool
  }{
    {4, false, []int{}, 0, []int{}, 0, false, nil, 0, nil, "", true},
    {4, false, []int{0, 1}, 2, []int{}, 2, false, nil, 2, 1, "", false},
    {4, false, []int{0, 1, 2}, 3, []int{3, 4}, 5, false, nil, 5, 4, "", false},
    {4, false, []int{0, 1, 2}, 3, []int{3, 4}, 5, true, 4, 4, 3, "", false},
    {4, false, []int{0}, 1, []int{}, 1, true, 0, 0, nil, "", true},
    {4, true, []int{0, 1}, 2, []int{3, 4}, 4, true, 4, 3, 3, "", false},
    {4, true, []int{0, 1, 2, 3}, 4, []int{4}, 4, false, nil, 4, 3, "bounded stack full: cannot push '4' onto the stack", false},
  }
  for i, test := range tests {
    s := NewStack(test.cap, test.bounded)
    var j, v int
    var val interface{}
    for j, v = range test.initVals {
      err := s.Push(v)
      if err != nil {
        t.Errorf("%d: unexpected error while pushing element %d: %q", i, j, err)
        goto Next
      }
    }
    if s.Size() != test.initSize {
      t.Errorf("%d: after initial push on stack, expected %d items, got %d", test.initSize, s.Size())
      continue
    }
    for j, v = range test.pushVals {
      err := s.Push(v)
      if err != nil && err.Error() != test.expectedErr {
        t.Errorf("%d: expected push of item %d onto stack to have an error of %q, got %q", i, j, test.expectedErr, err.Error())
        goto Next
      }
    }
    if s.Size() != test.pushedSize {
      t.Errorf("%d: after pushing more items onto the stack, expected the count to be %d, got %d", i, test.pushedSize, s.Size())
      continue
    }
    if test.pop {
      val = s.Pop()
      if val != test.popVal {
        t.Errorf("%d: expected the popped val to be %v, got %v", i, test.popVal, val)
      }
      if s.Size() != test.poppedSize {
        t.Errorf("%d: expected the count, after popping an item, to be %d, got %d", i, test.poppedSize, s.Size())
      }
    }
    val = s.Peek()
    if val != test.peekVal {
      t.Errorf("%d: expected the peeked item to be %v, got %v", i, test.peekVal, val)
    }
    if s.IsEmpty() != test.isEmpty {
      t.Errorf("%d: expected isEmpty to return %t, got %t", i, test.isEmpty, s.IsEmpty())
    }
Next:
  }
}
