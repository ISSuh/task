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
	prev atomic.Pointer[node]
	next atomic.Pointer[node]
}

func newNode(task *Task) *node {
	return &node{
		task: task,
	}
}

type TaskQueue struct {
	head atomic.Pointer[node]
	tail atomic.Pointer[node]
	size atomic.Int32
}

func NewTaskQueue() *TaskQueue {
	q := &TaskQueue{
		head: atomic.Pointer[node]{},
		tail: atomic.Pointer[node]{},
	}

	node := &node{}
	q.head.Store(node)
	q.tail.Store(node)
	return q
}

func (q *TaskQueue) Push(task *Task) {
	node := newNode(task)
	if q.Empty() {
		head := q.head.Load()
		tail := q.tail.Load()

		head.next.Store(node)
		tail.prev.Store(node)

		node.prev.Store(head)
		node.next.Store(tail)
	} else {
		tail := q.tail.Load()
		lastNode := tail.prev.Load()

		lastNode.next.CompareAndSwap(tail, node)
		tail.prev.CompareAndSwap(lastNode, node)

		node.prev.CompareAndSwap(nil, lastNode)
		node.next.CompareAndSwap(nil, tail)
	}

	q.size.Add(1)
}

func (q *TaskQueue) Pop() *Task {
	if q.Empty() {
		return nil
	}

	head := q.tail.Load()
	firstNode := head.next.Load()
	nextNode := firstNode.next.Load()

	head.next.CompareAndSwap(firstNode, nextNode)
	nextNode.prev.CompareAndSwap(firstNode, head)

	q.size.Add(-1)
	return firstNode.task
}

func (q *TaskQueue) Empty() bool {
	return q.size.Load() == 0
}

func (q *TaskQueue) Size() int {
	return int(q.size.Load())
}
