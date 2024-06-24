package main

import (
	"fmt"
	"reflect"
)

type task[T any] struct {
	Run T

	fn        reflect.Value
	fixedArgs []reflect.Value
}

func (t *task[T]) run(args []reflect.Value) []reflect.Value {
	allArgs := t.fixedArgs
	allArgs = append(allArgs, args...)
	for _, arg := range allArgs {
		fmt.Printf("arg : %+v\n", arg)
	}

	return t.fn.Call(allArgs)
}

// bind 함수 구현
func bind[T any](f interface{}, fixedArgs ...interface{}) task[T] {
	var funcDef T
	funcDefValue := reflect.ValueOf(funcDef)

	fnValue := reflect.ValueOf(f)
	fnType := fnValue.Type()

	fnOutsNum := fnValue.Type().NumOut()
	fnOuts := make([]reflect.Type, fnOutsNum)
	for i := 0; i < fnOutsNum; i++ {
		out := fnValue.Type().Out(i)
		fnOuts[i] = out
		fmt.Printf("outs : %s\n", fnOuts)
	}

	if fnType.Kind() != reflect.Func {
		panic("first argument must be a function")
	}

	if len(fixedArgs) >= fnType.NumIn() {
		panic("too many arguments to bind")
	}

	fixedArgsValue := make([]reflect.Value, len(fixedArgs))
	for i, arg := range fixedArgs {
		fixedArgsValue[i] = reflect.ValueOf(arg)
	}

	t := task[T]{
		fn:        fnValue,
		fixedArgs: fixedArgsValue,
	}

	fn := reflect.MakeFunc(funcDefValue.Type(), t.run).Interface().(T)

	// fn := reflect.MakeFunc(funcDefValue.Type(), func(args []reflect.Value) []reflect.Value {
	// 	allArgs := make([]reflect.Value, len(fixedArgs)+len(args))
	// 	for i, arg := range fixedArgs {
	// 		allArgs[i] = reflect.ValueOf(arg)
	// 	}
	// 	for i, arg := range args {
	// 		allArgs[len(fixedArgs)+i] = arg
	// 	}
	// 	return fnValue.Call(allArgs)
	// }).Interface().(T)

	t.Run = fn
	return t
}

func add(a, b int) (int, int) {
	return a, b
}

func main() {
	// add 함수의 첫 번째 인수를 5로 고정
	addFive := bind[func(int) (int, int)](add, 5)
	fmt.Println("before")
	fmt.Println(addFive.Run(3))

}
