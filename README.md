dq
=====
Dq got its name from dynamic queue, which was what this package originally implemented. The name remains the same.

DQ implements a queue, either bounded or unbounded.

DQ implements a stack, either bounded or unbounded.

## Queue
This implements a queue that can either be bounded or unbounded. The queue itself is an `[]interface{}`.

The design goals of this queue were:

* a queue that can grow as needed
* a queue that does not grow unnecessarily, i.e. if a certain percentage of the items in the queue has been dequeued, shift the remaining items in the queue forward so that new items can be enqueued without forcing a growth in the queue
* is safe for concurrent usage
* can act as a bounded queue

Reallocations are minimized by setting the initial capacity of the queue to a reasonable value for your use case.  Once a queue is grows, it does not shrink, even when the queue is emptied. Queue growth also results in any items in the queue being shifted forward in the slice to eliminate empty spaces in the front of the slice.

For unbounded queues, before growing the queue, the amount of empty space in the slice is checked and if it equals or exceeds the queue's shift percentage, instad of growing the slice, the items in the queue are shifted to the beginning of the slice.  By default, this shift percentage is set to 50%. This can be changed using the queue's `SetShiftPercent()` method.

For bounded queues, if the current queue length is equal to its capacity and there is an item to enqueue, the queue is checked to see if any elements have been dequeued.  If there is space at the beginning of the queue, all items are shifted forward, making room for the new item.  If the queue is full, an error is returned.

Operations supported:
```
* Enqueue
* Dequeue
* SetShiftPercent
* IsEmpty
* IsFull
* Size
* Reset
```
## Stack
This implements a stack that can either be bounded or unbounded. The stack itself is an `[]interface{}`.

A stack is created by calling `NewStack(size, bounded)`. The `size` is the initial capacity of the stack; `bounded` is a bool for whether or not this stack is bounded. If true, the stack will never be larger than its initial size. If false, the stack will grow as needed.

Operations supported:
```
* Push
* Pop
* Peek
* IsEmpty
* Size
```

For bounded queues, an error will occur on `Push()` operations if the queue is full.

`Pop()` and `Peek()` operations return both a value and a bool. If the stack is empty, an interface containing nil and false will be returned, otherwise the value and true will be returned.

## Usage
### Queue
A new queue can be obtained by either using either the `NewQ()` or `NewQueue()` functions; `NewQueue()` is an alias for `NewQ()`

Get an unbounded queue:

    q := dq.NewQ(256, false)

This returns a queue with an initial capacity of 256 items.

Get a bounded queue:

    q := dq.NewQueue(256, true)

This returns a bounded queue with a capacity of 256 items.

### Stack
A new stack can be obtained by using the `NewStack()` function.

Get an unbounded stack:

    s := dq.NewStack(256, false)

This returns a stack with an intial capacity of 256 items.

Get a bounded stack:

    s := dq.NewStack(256, false)

This returns a bounded stack with a capacity of 256 items.

## License
This code is licensed under the MIT license. For more information, please check the included LICENSE file.
