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
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var taskQueueCalledSlice []int
var taskQueueCalledMap map[int]bool

func resetTaskQueueCalled() {
	taskQueueCalledMap = map[int]bool{
		0: false,
		1: false,
		2: false,
	}

	taskQueueCalledSlice = []int{}
}

func a() {
	taskQueueCalledMap[0] = true
	taskQueueCalledSlice = append(taskQueueCalledSlice, 0)
}

func b() {
	taskQueueCalledMap[1] = true
	taskQueueCalledSlice = append(taskQueueCalledSlice, 1)
}

func c() {
	taskQueueCalledMap[2] = true
	taskQueueCalledSlice = append(taskQueueCalledSlice, 2)
}

func TestTaskQueue(t *testing.T) {

	taskA, err := Bind[TaskSigniture](a)
	require.NoError(t, err)

	taskB, err := Bind[TaskSigniture](b)
	require.NoError(t, err)

	taskC, err := Bind[TaskSigniture](c)
	require.NoError(t, err)

	t.Run("push", func(t *testing.T) {
		resetTaskQueueCalled()

		q := newTaskQueue()

		q.Push(NewTask(taskA))
		q.Push(NewTask(taskB))
		q.Push(NewTask(taskC))

		require.False(t, q.Empty())
		require.Equal(t, 3, q.Size())
	})

	t.Run("pop", func(t *testing.T) {
		resetTaskQueueCalled()

		q := newTaskQueue()

		q.Push(NewTask(taskA))
		q.Push(NewTask(taskB))
		q.Push(NewTask(taskC))

		require.False(t, q.Empty())
		require.Equal(t, 3, q.Size())

		task := q.Pop()
		task.Run()

		task = q.Pop()
		task.Run()

		task = q.Pop()
		task.Run()

		require.True(t, q.Empty())

		for i := 0; i < len(taskQueueCalledSlice); i++ {
			require.True(t, taskQueueCalledMap[0])
			require.Equal(t, i, taskQueueCalledSlice[i])
		}
	})

	t.Run("concurrent", func(t *testing.T) {
		t.Run("push", func(t *testing.T) {
			resetTaskQueueCalled()

			size := 100
			q := newTaskQueue()

			wg := sync.WaitGroup{}
			for i := 0; i < size; i++ {
				wg.Add(1)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()

					q.Push(NewTask(taskA))
				}(&wg)
			}

			wg.Wait()

			require.Equal(t, size, q.Size())
		})

		t.Run("pop", func(t *testing.T) {
			resetTaskQueueCalled()

			size := 100
			q := newTaskQueue()

			for i := 0; i < size; i++ {
				q.Push(NewTask(taskA))
			}

			wg := sync.WaitGroup{}
			for i := 0; i < size; i++ {
				wg.Add(1)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()

					task := q.Pop()
					require.NotNil(t, task)
				}(&wg)
			}

			wg.Wait()

			require.True(t, q.Empty())
		})
	})

	t.Run("clear", func(t *testing.T) {
		resetTaskQueueCalled()

		size := 100
		q := newTaskQueue()

		wg := sync.WaitGroup{}
		for i := 0; i < size; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				q.Push(NewTask(taskA))
			}(&wg)
		}

		wg.Wait()

		require.Equal(t, size, q.Size())

		q.Clear()
		require.True(t, q.Empty())
	})
}
