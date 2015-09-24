dq
=====
Dq got its name from dynamic queue, which was what this package originally implemented. The name remains the same.

DQ implements a queue, either bounded or unbounded.

DQ implements a stack, either bounded or unbounded.

All implementations are thread-safe.

## Queue
There are two queue implementations: unboundeed and bounded.  For each, the queue itself is an `[]interface{}`.

### Bounded queue
The bounded queue is implemented as a circular queue using a slice with a capacity that is one slot greater than the requested size. This allows for easy detection of whether or not the queue is full or empty.

If the queue is full, an error will be returned and the item will not be added to the queue. If, instead of an error, you wish to have the item replace the oldest item, then use the ring buffer.

During initial queue creation, all slots are initialized. This makes the intial queue request slower than just allocatin the memory for the queue but eliminates the need for additional logic in the queue to check whether or not the slot was already initialized, which is only useful the first time the queue is filled.

After queue creation, all item operations are done using the slice index.

Getting a circular queue:

    q := NewCircularQ(size)

Supported operations:
```
Enqueue(item)
Dequeue() (item bool)
Peek() (item bool)
IsEmpty() bool
IsFull() bool
Len() int
Cap() int
```

### Unbounded queue
The design goals of this queue were:

* a queue that can grow as needed
* a queue that does not grow unnecessarily, i.e. if a certain percentage of the items in the queue has been dequeued, shift the remaining items in the queue forward so that new items can be enqueued without forcing a growth in the queue
* is safe for concurrent usage
* can act as a bounded queue
* a queue from which memory can be reclaimed.

Reallocations are minimized by setting the initial capacity of the queue to a reasonable value for your use case.  Once a queue is grows, it does not shrink, even when the queue is emptied. Queue growth also results in any items in the queue being shifted forward in the slice to eliminate empty spaces in the front of the slice.

For unbounded queues, before growing the queue, the amount of empty space in the slice is checked and if it equals or exceeds the queue's shift percentage, instad of growing the slice, the items in the queue are shifted to the beginning of the slice.  By default, this shift percentage is set to 50%. This can be changed using the queue's `SetShiftPercent()` method.

For bounded queues, if the current queue length is equal to its capacity and there is an item to enqueue, the queue is checked to see if any elements have been dequeued.  If there is space at the beginning of the queue, all items are shifted forward, making room for the new item.  If the queue is full, an error is returned.

Getting a unbounded queue:

    q := queue.NewQueue(initialSize)

or

    q := queue.NewQ(initialSize)

Operations supported:
```
Enqueue(item) error
Dequeue() (item, bool)
Peek() (item, bool)
SetShiftPercent(int)
IsEmpty() bool
IsFull() bool
Len() int
Cap() int
Reset()
```
## Stack
This implements a stack that can either be bounded or unbounded. The stack itself is an `[]interface{}`.

A stack is created by calling `NewStack(size, bounded)`. The `size` is the initial capacity of the stack; `bounded` is a bool for whether or not this stack is bounded. If true, the stack will never be larger than its initial size. If false, the stack will grow as needed.

Getting a bounded stack with 256 slots:

    s := stack.NewStack(256, true)

Getting an unboudned stack with an initial capacity of 256 slots:

    s := stack.NewStack(256, false)

Operations supported:
```
Push(item) error
Pop() (interface{}, bool)
Peek() (interface{}, bool)
IsEmpty() bool
IsFull() bool
Len() int
Cap() int
Reset()
```

For bounded queues, an error will occur on `Push()` operations if the queue is full.

`Pop()` and `Peek()` operations return both a value and a bool.  If the stack is empty, an interface containing nil and false will be returned, otherwise the value and true will be returned.

## Buffer
Buffer implements a ring buffer using a `[]interface{}`.  This is a wrapper to quueue.Circular.  See that for more infomration.

When full, instead of erroring, like the ciruclar queue, the ring buffer evicts the oldest item in the buffer and enqueues the new item at the back of the buffer.

Getting a ring buffer with 256 slots:

    ring := buffer.Ring(256)

Operations supported:
```
Enqueue(interface{}) error
Dequeue() (interface{}, bool)
Peek() (interface{}, bool)
IsEmpty() bool
IsFull() bool
Len() int
Cap() int
Reset()

## License
This code is licensed under the MIT license. For more information, please check the included LICENSE file.
