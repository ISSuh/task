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
	"context"
	"fmt"
	"sync"
	"time"
)

type taskWorker struct {
	id     int
	queue  *taskQueue
	notify <-chan struct{}
}

func newTaskWorker(id int, q *taskQueue, notify <-chan struct{}) *taskWorker {
	return &taskWorker{
		id:     id,
		queue:  q,
		notify: notify,
	}
}

func (w *taskWorker) work(c context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-c.Done():
			fmt.Printf("[%d] context canceled\n", w.id)
			return
		case <-w.notify:
			w.workInternal()
		}
	}
}

func (w *taskWorker) workInternal() {
	task := w.queue.Pop()
	if task == nil {
		return
	}

	if time.Now().After(task.timestamp) {
		task.callback.Run()
	} else {
		w.queue.Push(task)
	}
}
