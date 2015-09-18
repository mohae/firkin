dq
=====
Dq is a dynamic queue that supports growth, compaction, and is concurrency safe. The size of the queue can be constrained by setting its `maxCapacity` to a value > 0. This is set at queue creation time via the `New()` function.

The queue itself is `[]interface{}`.

The design goals of this queue were:

* a queue that can grow as needed
* a queue that does not grow unnecessarily, i.e. if a certain percentage of the items in the queue has been dequeued, shift the remaining items in the queue forward so that new items can be enqueued without forcing a growth in the queue
* is safe for concurrent usage

Reallocations are minimized by setting the initial capacity of the queue to a reasonable value for your use case. Any queue growth that occurs after queue creation follows the algorithm in Go's growSlice(). Once a queue is grown, it is not shrunk back, even when the queue is emptied. Any queue growth also results in any items in the queue being shifted forward in the slice to eliminate empty spaces in the front of the slice.

Before growing the queue, the amount of empty space in the slice is checked and if it exceeds the queue's shift percentage, the items in the queue are shifted to the beginning of the slice, avoiding the allocations due to growth. This shift percentage defaults to 50%; it can be set using the queue's SetShiftPercent() method.

## Usage
Go get:

    go get github.com/mohae/dq

Import:

    import github.com/mohae/dq

Get a queue:

    q := dq.New(256, 0)

This returns a queue with an initial capacity of 256 items and without a maximum capacity.

## License
This code is licensed under the MIT license. For more information, please check the included LICENSE file.
