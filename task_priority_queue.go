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
	"container/heap"
	"sync"
)

type taskPriorityQueue []*Task

var mutex sync.Mutex

func newTaskPriorityQueue() taskPriorityQueue {
	q := make(taskPriorityQueue, 0)
	heap.Init(&q)
	return q
}

func (q taskPriorityQueue) Len() int {
	return len(q)
}

func (q taskPriorityQueue) Less(i, j int) bool {
	return q[i].isBefore(q[j])
}

func (q taskPriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *taskPriorityQueue) Push(x interface{}) {
	item := x.(*Task)
	*q = append(*q, item)
}

func (q *taskPriorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[0 : n-1]
	return item
}

func (q *taskPriorityQueue) PushTask(task *Task) {
	mutex.Lock()
	defer mutex.Unlock()

	heap.Push(q, task)
}

func (q *taskPriorityQueue) PopTask() *Task {
	mutex.Lock()
	defer mutex.Unlock()

	return heap.Pop(q).(*Task)
}
