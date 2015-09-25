firkin
======
Firkin is a small package of containers.

Firkin implements a queue, either bounded or unbounded, as a ring queue.

Firkin implements a stack, either bounded or unbounded.

Firkin implements a ring buffer.

All implementations are thread-safe.

## Queue
There are two queue implementations: unboundeed and bounded.  For each, the queue itself is an `[]interface{}`.  All queue methods are thread-safe.

Queues can be resized by using the resize method: `Resize(newSize)`.  The newSize must be equal to or larger than both 1.25 * the number of elements in the queue or the intial queue capacity, whichever is larger.  Use `0` as the newSize if you want either 1.25 * the number of elements in the queue or the intial queue capacity used.  Memory may be reclaimed during a Resize operation.  Memory may also be allocated during a Resize operation.  The queue's new size is returned.

Queues can be reset. Queue reset causes all items in the queue to be lost. A reset will not reclaim memory.

Supported operations:
```
Enqueue(item)
Dequeue() (item bool)
Peek() (item bool)
IsEmpty() bool
IsFull() bool
Len() int
Cap() int
Reset()
Resize(int) int
  ```

### Circular (Bounded) queue
The bounded queue is implemented as a circular queue using a slice with a capacity that is one slot greater than the requested size. This allows for easy detection of whether or not the queue is full or empty.

If the queue is full, an error will be returned and the item will not be added to the queue. If, instead of an error, you wish to have the item replace the oldest item, then use the ring buffer.

During initial queue creation, all slots are initialized. This makes the intial queue request slower than just allocatin the memory for the queue but eliminates the need for additional logic in the queue to check whether or not the slot was already initialized, which is only useful the first time the queue is filled.

After queue creation, all item operations are done using the slice index.

When resizing a circular queue, the new queue slots are zero'd, this is an `O(n)` process where n is the new queue capacity. Any items in the queue are copied to the front of the new queue.

Getting a circular queue:

    q := NewCircularQ(size)

For bounded queues, if the current queue length is equal to its capacity and there is an item to enqueue, the queue is checked to see if any elements have been dequeued.  If there is space at the beginning of the queue, all items are shifted forward, making room for the new item.  If the queue is full, an error is returned.

Bounded queues can be resized using the `Resize(size)` method.  Bounded queues do not automatically resize.  Resize operations allow the queue to grow or shrink. For a buffer to successfully shrink, there most be less items left in the buffer than the new buffer size.  During resize operations, any items in the buffer will be copied to a tmp buffer and then recopied to the resized queue.

### Unbounded queue
The design goals of this queue were:

* a queue that can grow as needed
* a queue that does not grow unnecessarily, i.e. if a certain percentage of the items in the queue has been dequeued, shift the remaining items in the queue forward so that new items can be enqueued without forcing a growth in the queue
* a queue from which memory can be reclaimed.

Reallocations are minimized by setting the initial capacity of the queue to a reasonable value for your use case.  Once a queue is grows, it does not shrink, even when the queue is emptied. Queue growth also results in any items in the queue being shifted forward in the slice to eliminate empty spaces in the front of the slice.

For unbounded queues, before growing the queue, the amount of empty space in the slice is checked and if it equals or exceeds the queue's shift percentage, instead of growing the slice, the items in the queue are shifted to the beginning of the slice.  By default, this shift percentage is set to 50%. This can be changed using the queue's `SetShiftPercent()` method.

When a queue is resized, all of the elements in the existing queue, if there are any, are copied to the front of the new queue. This is an `O(n)` operation.

Getting a unbounded queue:

    q := queue.NewQueue(initialSize)

or

    q := queue.NewQ(initialSize)

Additional supported operations:
```
SetShiftPercent(int)
```

## Stack
This implements a stack that can either be bounded or unbounded. The stack itself is an `[]interface{}`.

A stack is created by calling `NewStack(size, bounded)`. The `size` is the initial capacity of the stack; `bounded` is a bool for whether or not this stack is bounded. If true, the stack will never be larger than its initial size. If false, the stack will grow as needed.

Stacks can be resized by using the resize method: `Resize(newSize)`.  The newSize must be equal to or larger than both 1.25 * the number of elements in the stack or the intial queue capacity, whichever is larger.  Use `0` as the newSize if you want either 1.25 * the number of elements in the stack or the intial stack capacity used.  The stack's new size is returned.

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
Resize(int) int
```

### Bounded Stack
For bounded queues, an error will occur on `Push()` operations if the queue is full.

Getting a bounded stack with 256 slots:

    s := stack.NewStack(256, true)

### Unbounded Stack
Getting an unbound stack with an initial capacity of 256 slots:

    s := stack.NewStack(256, false)

## Buffer
Buffer implements a ring buffer using a `[]interface{}`.  This is a wrapper to queue.Circular.  See that for more infomration.

When full, instead of creating an error, like the circular queue, the ring buffer evicts the oldest item in the buffer and enqueues the new item at the back of the buffer.

Getting a ring buffer with 256 slots:

    ring := buffer.Ring(256)

## License
This code is licensed under the MIT license. For more information, please check the included LICENSE file.
