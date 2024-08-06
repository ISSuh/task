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

import "time"

type TaskSigniture func()

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type Task struct {
	_ noCopy

	timestamp time.Time
	callback  *callback[TaskSigniture]
}

func NewTask(callback *callback[TaskSigniture]) *Task {
	return &Task{
		timestamp: time.Now(),
		callback:  callback,
	}
}

func NewDelayTask(delay time.Duration, callback *callback[TaskSigniture]) *Task {
	return &Task{
		timestamp: time.Now().Add(delay),
		callback:  callback,
	}
}

func (t *Task) Run() {
	t.callback.Run()
}

func (t *Task) isBefore(b *Task) bool {
	return t.timestamp.Before(b.timestamp)
}

func (t *Task) TimeStamp() time.Time {
	return t.timestamp
}
