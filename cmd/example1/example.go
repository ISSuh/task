package main

import (
	"fmt"
	"reflect"
)

// type functor func(args ...any)

type task[T any] struct {
	callback reflect.Value
	fixArgs  []reflect.Value
}

// func (t *task[T]) Run(args ...any) {
// 	// fixArgs := t.fixArgs
// 	allArgs := []reflect.Value{}
// 	for _, arg := range args {
// 		argValue := reflect.ValueOf(arg)
// 		allArgs = append(allArgs, argValue)
// 	}

// 	result := t.callback.Call(allArgs)
// 	for _, r := range result {
// 		fmt.Printf("%d\n", r.Interface().(int))
// 	}

// }

// bind 함수 구현
func bind[T any](f interface{}, fixedArgs ...interface{}) task[T] {
	var funcDef T
	funcDefValue := reflect.ValueOf(funcDef)
	fnValue := reflect.ValueOf(f)

	if fnValue.Kind() != reflect.Func {
		panic("first argument must be a function")
	}

	fixedArgsValue := make([]reflect.Value, len(fixedArgs))
	for i, arg := range fixedArgs {
		fixedArgsValue[i] = reflect.ValueOf(arg)
	}

	fn := reflect.MakeFunc(funcDefValue.Type(), func(args []reflect.Value) []reflect.Value {
		allArgs := make([]reflect.Value, len(fixedArgs)+len(args))
		for i, arg := range fixedArgs {
			allArgs[i] = reflect.ValueOf(arg)
		}
		for i, arg := range args {
			allArgs[len(fixedArgs)+i] = arg
		}
		return fnValue.Call(allArgs)
	})

	task := task[T]{
		callback: fn,
		fixArgs:  fixedArgsValue,
	}
	return task
}

func add(a, b int) int {
	return a + b
}

// func mul(a, b, c int) float32 {
// 	return float32(a*b*c) * 1.0
// }

// type Adder int

// func (q Adder) add(a, b int) float32 {
// 	return int(q) + a + b
// }

func TesmpArgs(args ...any) int {
	return 0
}

func main() {
	a := bind[func(int) int](add, 5)
	a.Run(3)
	fmt.Println()

	// m := bind[func(int, int) float32](mul, 1)
	// fmt.Println(m.callback(2, 3))

	// addera := Adder(1)
	// addaa := bind[func(int) int](addera.add, 2)
	// fmt.Println(addaa(3)) // 출력: 8
}
