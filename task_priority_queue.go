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

import "container/heap"

type item struct {
	task  *Task
	index int
}

type taskPriorityQueue []*item

func newTaskPriorityQueue() taskPriorityQueue {
	q := make(taskPriorityQueue, 0)
	heap.Init(&q)
	return q
}

func (q taskPriorityQueue) Len() int {
	return len(q)
}

func (q taskPriorityQueue) Less(i, j int) bool {
	return q[i].task.isBefore(q[j].task)
}

func (q taskPriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *taskPriorityQueue) Push(x interface{}) {
	item := x.(*item)
	// item.index = len(*q)
	*q = append(*q, item)
}

func (q *taskPriorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	// item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}

// 우선순위 큐에 아이템 추가
func (q *taskPriorityQueue) PushItem(item *item) {
	heap.Push(q, item)
}

// 우선순위 큐에서 아이템 추출 (우선순위가 가장 높은 아이템)
func (q *taskPriorityQueue) PopItem() *item {
	return heap.Pop(q).(*item)
}
