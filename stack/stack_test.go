package stack

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
    popOk bool
    poppedSize int
    peekVal interface{}
    peekOk bool
    expectedErr string
    isEmpty bool
  }{
    {4, false, []int{}, 0, []int{}, 0, false, nil, false, 0, nil, false, "", true},
    {4, false, []int{0, 1}, 2, []int{}, 2, false, nil, false, 2, 1, true, "", false},
    {4, false, []int{0, 1, 2}, 3, []int{3, 4}, 5, false, nil, false, 5, 4, true, "", false},
    {4, false, []int{0, 1, 2}, 3, []int{3, 4}, 5, true, 4, true, 4, 3, true, "", false},
    {4, false, []int{0}, 1, []int{}, 1, true, 0, true, 0, nil, false, "", true},
    {4, true, []int{0, 1}, 2, []int{3, 4}, 4, true, 4, true, 3, 3, true, "", false},
    {4, true, []int{0, 1, 2, 3}, 4, []int{4}, 4, false, nil, false, 4, 3, true, "bounded stack full: cannot push '4' onto the stack", false},
  }
  for i, test := range tests {
    s := NewStack(test.cap, test.bounded)
    var j, v int
    var val interface{}
    var ok bool
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
      val, ok = s.Pop()
      if ok != test.popOk {
        t.Errorf("%d: expected pop ok to be %t, got %t", i, test.popOk, ok)
      }
      if ok {
        if val != test.popVal {
          t.Errorf("%d: expected the popped val to be %v, got %v", i, test.popVal, val)
        }
      }
      if s.Size() != test.poppedSize {
        t.Errorf("%d: expected the count, after popping an item, to be %d, got %d", i, test.poppedSize, s.Size())
      }
    }
    val, ok = s.Peek()
    if ok != test.peekOk {
      t.Errorf("%d: expected peep ok to be %t, got %t", i, test.peekOk, ok)
    }
    if ok {
      if val != test.peekVal {
        t.Errorf("%d: expected the peeked item to be %v, got %v", i, test.peekVal, val)
      }
    }
    if s.IsEmpty() != test.isEmpty {
      t.Errorf("%d: expected isEmpty to return %t, got %t", i, test.isEmpty, s.IsEmpty())
    }
Next:
  }
}
