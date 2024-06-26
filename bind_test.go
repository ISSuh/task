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
	"testing"

	"github.com/stretchr/testify/require"
)

type testFunctions struct {
	Called map[int]bool
}

func (t *testFunctions) reset() {
	t.Called = map[int]bool{}
}

func (t *testFunctions) func1() {
	t.Called[1] = true
}

func (t *testFunctions) func2(a int) {
	t.Called[2] = true
}

func (t *testFunctions) func3(a int) int {
	t.Called[3] = true
	return a
}

func (t *testFunctions) func4(a int, b float32) (int, float32) {
	t.Called[4] = true
	return a, b
}

func TestBind(t *testing.T) {
	funcs := testFunctions{}

	t.Run("basic function", func(t *testing.T) {
		t.Run("binding parameter", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func()
			callback, err := Bind[FuntionDefinition](funcs.func1)
			require.NoError(t, err)

			callback.Run()

			require.NotNil(t, callback)
			require.True(t, funcs.Called[1])
		})

		t.Run("invalid parameter binding", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func(int)
			callback, err := Bind[FuntionDefinition](funcs.func1)
			require.Error(t, err)
			require.Nil(t, callback)
		})

		t.Run("invalid return binding", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func() int
			callback, err := Bind[FuntionDefinition](funcs.func1)
			require.Error(t, err)
			require.Nil(t, callback)
		})

		t.Run("invalid number of fixed parameter", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func()
			callback, err := Bind[FuntionDefinition](funcs.func1, 1)
			require.Error(t, err)
			require.Nil(t, callback)
		})

		t.Run("invalid generic type", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition int
			callback, err := Bind[FuntionDefinition](funcs.func1, 1)
			require.Error(t, err)
			require.Nil(t, callback)
		})
	})

	t.Run("one parameter function", func(t *testing.T) {
		t.Run("binding parameter", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func()
			callback, err := Bind[FuntionDefinition](funcs.func2, 5)
			require.NoError(t, err)

			callback.Run()

			require.NotNil(t, callback)
			require.True(t, funcs.Called[2])
		})

		t.Run("no binding parameter", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func(int)
			callback, err := Bind[FuntionDefinition](funcs.func2)
			require.NoError(t, err)

			callback.Run(1)

			require.NotNil(t, callback)
			require.True(t, funcs.Called[2])
		})
	})

	t.Run("one parameter, one return function", func(t *testing.T) {
		t.Run("binding parameter", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func() int
			callback, err := Bind[FuntionDefinition](funcs.func3, 5)
			require.NoError(t, err)

			callback.Run()

			require.NotNil(t, callback)
			require.True(t, funcs.Called[3])
		})

		t.Run("no binding parameter", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func(int) int
			callback, err := Bind[FuntionDefinition](funcs.func3)
			require.NoError(t, err)

			callback.Run(1)

			require.NotNil(t, callback)
			require.True(t, funcs.Called[3])
		})

		t.Run("empty return definition", func(t *testing.T) {
			funcs.reset()

			type FuntionDefinition func(int)
			callback, err := Bind[FuntionDefinition](funcs.func3)
			require.Error(t, err)
			require.Nil(t, callback)
		})
	})
}
