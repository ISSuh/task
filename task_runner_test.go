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
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var runnerTestACalledCount atomic.Int32
var runnerTestSquence string

func resetRunnerTest() {
	runnerTestACalledCount.Store(0)
	runnerTestSquence = ""
}

func runnerTestA() {
	fmt.Println("A")
	runnerTestACalledCount.Add(1)
	runnerTestSquence += "A"
}

func runnerTestB() {
	fmt.Println("B")
	runnerTestSquence += "B"
}

func runnerTestC() {
	fmt.Println("C")
	runnerTestSquence += "C"
}

func TestTaskRunner(t *testing.T) {

	task, err := Bind[TaskSigniture](runnerTestA)
	require.NoError(t, err)

	t.Run("cancel", func(t *testing.T) {
		resetRunnerTest()

		c, cancel := context.WithCancel(context.Background())
		r := NewTaskRunner(5)

		go func() {
			r.PostTask(NewTask(task))
			time.Sleep(1 * time.Second)
			cancel()
		}()

		r.RunLoop(c)

		require.Equal(t, 1, int(runnerTestACalledCount.Load()))
	})

	t.Run("post task", func(t *testing.T) {
		resetRunnerTest()

		size := 10
		c, cancel := context.WithCancel(context.Background())
		r := NewTaskRunner(5)

		wg := sync.WaitGroup{}
		for i := 0; i < size; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				r.PostTask(NewTask(task))
			}(&wg)
		}

		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()

		r.RunLoop(c)

		require.Equal(t, size, int(runnerTestACalledCount.Load()))
	})

	t.Run("post delay task", func(t *testing.T) {
		resetRunnerTest()

		c, cancel := context.WithCancel(context.Background())
		r := NewTaskRunner(5)

		taskA, err := Bind[TaskSigniture](runnerTestA)
		require.NoError(t, err)

		taskB, err := Bind[TaskSigniture](runnerTestB)
		require.NoError(t, err)

		taskC, err := Bind[TaskSigniture](runnerTestC)
		require.NoError(t, err)

		r.PostTask(NewDelayTask(100*time.Millisecond, taskA))
		r.PostTask(NewDelayTask(50*time.Millisecond, taskB))
		r.PostTask(NewTask(taskC))

		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()

		r.RunLoop(c)

		expectSequence := "CBA"
		require.Equal(t, expectSequence, runnerTestSquence)

		cancel()
	})
}
