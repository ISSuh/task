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

package main

import (
	"fmt"
	"reflect"
)

// 바인딩할 함수를 반환
func bind(f interface{}, fixedArgs ...interface{}) interface{} {
	fnValue := reflect.ValueOf(f)
	if fnValue.Kind() != reflect.Func {
		panic("first argument must be a function")
	}
	return reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) []reflect.Value {
		allArgs := make([]reflect.Value, len(fixedArgs)+len(args))
		for i, arg := range fixedArgs {
			allArgs[i] = reflect.ValueOf(arg)
		}
		for i, arg := range args {
			allArgs[len(fixedArgs)+i] = arg
		}
		return fnValue.Call(allArgs)
	}).Interface()
}

func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

func main() {
	// add 함수의 첫 번째 인수를 5로 고정
	addFive := bind(add, 5).(func(int) int)
	fmt.Println(addFive(3)) // 출력: 8

	// multiply 함수의 첫 번째 인수를 2로 고정
	multiplyByTwo := bind(multiply, 2).(func(int) int)
	fmt.Println(multiplyByTwo(4)) // 출력: 8

	// multiply 함수의 첫 번째 인수를 3으로 고정
	multiplyByThree := bind(multiply, 3).(func(int) int)
	fmt.Println(multiplyByThree(4)) // 출력: 12
}
