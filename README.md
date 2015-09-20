dq
=====
Dq got its name from dynamic queue, which was what this package originally implemented. The queue within this package can either function as a dynamic queue or a bounded queue.

## Queue
This implements a queue that can either be bounded, or unbounded. The queue itself is an `[]interface{}`.

The design goals of this queue were:

* a queue that can grow as needed
* a queue that does not grow unnecessarily, i.e. if a certain percentage of the items in the queue has been dequeued, shift the remaining items in the queue forward so that new items can be enqueued without forcing a growth in the queue
* is safe for concurrent usage
* can act as a bounded queue

Reallocations are minimized by setting the initial capacity of the queue to a reasonable value for your use case.  Once a queue is grows, it does not shrink, even when the queue is emptied. Queue growth also results in any items in the queue being shifted forward in the slice to eliminate empty spaces in the front of the slice.

For unbounded queues, before growing the queue, the amount of empty space in the slice is checked and if it equals or exceeds the queue's shift percentage, instad of growing the slice, the items in the queue are shifted to the beginning of the slice.  By default, this shift percentage is set to 50%. This can be changed using the queue's `SetShiftPercent()` method.

For bounded queues, if the current queue length is equal to its capacity and there is an item to enqueue, the queue is checked to see if any elements have been dequeued.  If there is space at the beginning of the queue, all items are shifted forward, making room for the new item.  If the queue is full, an error is returned.

### Usage
A new queue can be obtained by either using either the `NewQ()` or `NewQueue()` functions; `NewQueue()` is an alias for `NewQ()`
Get an unbounded queue:

    q := dq.NewQ(256, 0)

This returns a queue with an initial capacity of 256 items and without a maximum capacity.

Get a bounded queue:

    q := dq.NewQ(64, 256)

This returns a queue with an initial capacity of 64 items that can grow to a queue with a maximum size of 256 items.

## License
This code is licensed under the MIT license. For more information, please check the included LICENSE file.
