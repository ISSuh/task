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

	itemA := item{
		task:  NewDelayTask(5000*time.Millisecond, taskA),
		index: 0,
	}

	itemB := item{
		task:  NewDelayTask(2000*time.Millisecond, taskB),
		index: 1,
	}

	itemC := item{
		task:  NewDelayTask(1000*time.Millisecond, taskC),
		index: 2,
	}

	q.PushItem(&itemA)
	q.PushItem(&itemB)
	q.PushItem(&itemC)

	AA := q.PopItem()
	BB := q.PopItem()
	CC := q.PopItem()

	AA.task.Run()
	BB.task.Run()
	CC.task.Run()

	require.Equal(t, "CBA", runnerTestSquenceTmp)
}
