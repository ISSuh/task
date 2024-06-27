package main

import (
	"fmt"
	"reflect"
)

// type functor func(args ...any)

type task struct {
	Callback interface{}
	fixArgs  []reflect.Value
}

// bind 함수 구현
func bind(d interface{}, f interface{}, fixedArgs ...interface{}) task {
	funcDefValue := reflect.ValueOf(d).Elem()
	fnValue := reflect.ValueOf(f)

	if fnValue.Kind() != reflect.Func {
		panic("first argument must be a function")
	}

	fixedArgsValue := make([]reflect.Value, len(fixedArgs))
	for i, arg := range fixedArgs {
		fixedArgsValue[i] = reflect.ValueOf(arg)
	}

	fn := reflect.MakeFunc(
		funcDefValue.Type(),
		func(args []reflect.Value) []reflect.Value {
			allArgs := make([]reflect.Value, len(fixedArgs)+len(args))
			for i, arg := range fixedArgs {
				allArgs[i] = reflect.ValueOf(arg)
			}
			for i, arg := range args {
				allArgs[len(fixedArgs)+i] = arg
			}
			return fnValue.Call(allArgs)
		}).Interface()

	task := task{
		Callback: fn,
		fixArgs:  fixedArgsValue,
	}
	return task
}

func add(a, b int) int {
	return a + b
}

func main() {
	var Def func(int) int
	a := bind(&Def, add, 5)
	a.Callback.(func(int) int)(1)

	fmt.Println()
}
