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
)

type f interface{}

type funcor interface {
	f
}

func Bind[T funcor](f interface{}, fixedArgs ...interface{}) T {
	var functor T
	definitionValue := reflect.ValueOf(definition)

	funcValue := reflect.ValueOf(f)
	if funcValue.Kind() != reflect.Func {
		panic("first argument must be a function")
	}

	fn := reflect.MakeFunc(definitionValue.Type(), func(args []reflect.Value) []reflect.Value {
		allArgs := make([]reflect.Value, len(fixedArgs)+len(args))
		for i, arg := range fixedArgs {
			allArgs[i] = reflect.ValueOf(arg)
		}
		for i, arg := range args {
			allArgs[len(fixedArgs)+i] = arg
		}
		return funcValue.Call(allArgs)
	})
	return fn.Interface().(T)
}
