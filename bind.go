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
	"reflect"
)

func Bind[FuncDefinition any](f interface{}, fixedArgs ...interface{}) (*Task[FuncDefinition], error) {
	var definition FuncDefinition
	definitionValue := reflect.ValueOf(definition)
	definitionType := definitionValue.Type()

	funcValue := reflect.ValueOf(f)
	funcType := funcValue.Type()

	if funcType.Kind() != reflect.Func {
		return nil, fmt.Errorf("first argument must be a function")
	}

	predictedFixedArgsNum := funcType.NumIn() - definitionType.NumIn()
	if predictedFixedArgsNum != len(fixedArgs) {
		if len(fixedArgs) < predictedFixedArgsNum {
			return nil, fmt.Errorf("too few arguments to bind")
		} else {
			return nil, fmt.Errorf("too many arguments to bind")
		}
	}

	if err := compareReturnTypes(funcType, definitionType); err != nil {
		return nil, err
	}

	// fnOutsNum := funcType.NumOut()
	// fnOuts := make([]reflect.Type, fnOutsNum)
	// for i := 0; i < fnOutsNum; i++ {
	// 	out := funcType.Out(i)
	// 	fnOuts[i] = out
	// 	fmt.Printf("outs : %s\n", fnOuts)
	// }

	fixedArgsValues := make([]reflect.Value, len(fixedArgs))
	fixedArgsTypes := make([]reflect.Type, len(fixedArgs))
	funcArgsValues := make([]reflect.Type, len(fixedArgs))
	for i, arg := range fixedArgs {
		argValue := reflect.ValueOf(arg)

		fixedArgsValues[i] = argValue
		fixedArgsTypes[i] = argValue.Type()
		funcArgsValues[i] = funcType.In(i)
	}

	if err := compareArgumentsType(fixedArgsTypes, funcArgsValues, len(fixedArgsTypes)); err != nil {
		return nil, err
	}

	t := &Task[FuncDefinition]{
		function:  funcValue,
		fixedArgs: fixedArgsValues,
	}

	t.Run = reflect.MakeFunc(definitionValue.Type(), t.wrapper).Interface().(FuncDefinition)
	return t, nil
}

func isTypeEqual(a reflect.Type, b reflect.Type) bool {
	return a.Kind() == b.Kind()
}

func compareArgumentsType(a []reflect.Type, b []reflect.Type, length int) error {
	if len(a) < length || len(b) < length {
		return fmt.Errorf("invalid length")
	}

	for i := 0; i < length; i++ {
		if !isTypeEqual(a[i], b[i]) {
			return fmt.Errorf("not matched argument")
		}
	}
	return nil
}

func compareReturnTypes(a reflect.Type, b reflect.Type) error {
	if a.NumOut() != b.NumOut() {
		return fmt.Errorf("not matched number of returns")
	}

	for i := 0; i < a.NumOut(); i++ {
		aOutType := a.Out(i)
		bOutType := b.Out(i)
		if !isTypeEqual(aOutType, bOutType) {
			return fmt.Errorf("not matched %d's return type", i)
		}
	}
	return nil
}
