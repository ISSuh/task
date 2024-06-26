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
	"reflect"
	"sync/atomic"
)

type Callback[T any] struct {
	Run       T
	function  reflect.Value
	fixedArgs []reflect.Value
}

func (t *Callback[T]) helper(args []reflect.Value) []reflect.Value {
	allArgs := t.fixedArgs
	allArgs = append(allArgs, args...)
	return t.function.Call(allArgs)
}

type OnceCallback[T any] struct {
	*Callback[T]

	called atomic.Bool
}

func (t *OnceCallback[T]) helper(args []reflect.Value) []reflect.Value {
	if t.called.Load() {
		panic("already called")
	}

	t.called.Store(true)
	return t.Callback.helper(args)
}
