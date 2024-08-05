// MIT License

// Copyright (c) 2024 ISSuh

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package worker

import (
	"sync/atomic"
)

type node struct {
	task *Task
	next atomic.Pointer[node]
}

func newNode(task *Task) *node {
	return &node{
		task: task,
	}
}

type taskQueue struct {
	head atomic.Pointer[node]
	tail atomic.Pointer[node]
	size atomic.Int32
}

func newTaskQueue() *taskQueue {
	q := taskQueue{}

	q.Clear()
	return &q
}

func (q *taskQueue) Push(task *Task) {
	node := newNode(task)
	for {
		tail := q.tail.Load()
		next := tail.next.Load()
		if next == nil {
			if tail.next.CompareAndSwap(next, node) {
				q.tail.CompareAndSwap(tail, node)
				q.size.Add(1)
				return
			}
		} else {
			q.tail.CompareAndSwap(tail, next)
		}
	}
}

func (q *taskQueue) Pop() *Task {
	if q.Empty() {
		return nil
	}

	for {
		tail := q.tail.Load()
		head := q.head.Load()
		next := head.next.Load()
		if head != tail {
			if next == nil {
				return nil
			}

			task := next.task
			if q.head.CompareAndSwap(head, next) {
				q.size.Add(-1)
				return task
			}
		} else {
			if next == nil {
				return nil
			}

			q.tail.CompareAndSwap(tail, next)
		}
	}
}

func (q *taskQueue) Clear() {
	node := &node{}
	q.head.Store(node)
	q.tail.Store(node)
	q.size.Store(0)
}

func (q *taskQueue) Empty() bool {
	return q.size.Load() == 0
}

func (q *taskQueue) Size() int {
	return int(q.size.Load())
}
