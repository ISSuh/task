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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var runnerTestSquenceTmp string

func runnerTestAA() {
	fmt.Println("A")
	runnerTestSquenceTmp += "A"
}

func runnerTestBB() {
	fmt.Println("B")
	runnerTestSquenceTmp += "B"
}

func runnerTestCC() {
	fmt.Println("C")
	runnerTestSquenceTmp += "C"
}

func Test_taskPriorityQueue(t *testing.T) {
	q := newTaskPriorityQueue()

	taskA, err := Bind[TaskSigniture](runnerTestAA)
	require.NoError(t, err)

	taskB, err := Bind[TaskSigniture](runnerTestBB)
	require.NoError(t, err)

	taskC, err := Bind[TaskSigniture](runnerTestCC)
	require.NoError(t, err)

	itemA := NewDelayTask(1000*time.Millisecond, taskA)
	itemB := NewDelayTask(500*time.Millisecond, taskB)
	itemC := NewDelayTask(100*time.Millisecond, taskC)

	q.PushTask(itemA)
	q.PushTask(itemB)
	q.PushTask(itemC)

	AA := q.PopTask()
	BB := q.PopTask()
	CC := q.PopTask()

	AA.Run()
	BB.Run()
	CC.Run()

	require.Equal(t, "CBA", runnerTestSquenceTmp)
}
