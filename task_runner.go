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
)

type TaskRunner struct {
	queue   taskQueue
	workers []*taskWorker

	workerNum    int
	cancelWorker context.CancelFunc
	notify       chan struct{}
	wg           sync.WaitGroup
}

func NewTaskRunner(workerNum int) *TaskRunner {
	r := &TaskRunner{
		queue:     newTaskQueue(),
		workerNum: workerNum,
		workers:   make([]*taskWorker, workerNum),
		notify:    make(chan struct{}, workerNum),
		wg:        sync.WaitGroup{},
	}

	c, cancel := context.WithCancel(context.Background())
	r.cancelWorker = cancel

	for i := 0; i < workerNum; i++ {
		r.wg.Add(1)
		r.workers[i] = newTaskWorker(i, &r.queue, r.notify)
		go r.workers[i].work(c, &r.wg)
	}
	return r
}

func (r *TaskRunner) PostTask(task *Task) {
	r.queue.Push(task)
	r.notify <- struct{}{}
}

func (r *TaskRunner) RunLoop(c context.Context) {
	<-c.Done()

	fmt.Printf("run lopp finished\n")
	r.cancelWorker()

	r.wg.Wait()
}
