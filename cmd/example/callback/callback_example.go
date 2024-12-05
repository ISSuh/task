package main

import (
	"fmt"

	worker "github.com/ISSuh/worker"
)

func add(a, b int) int {
	return a + b
}

func main() {
	// Bind function and return callback
	callback1, err := worker.Bind[func(a, b int) int](add)
	if err != nil {
		panic(err)
	}

	// Run callback with parameter
	res1 := callback1.Run(1, 2)
	fmt.Printf("callback1 result: %d\n", res1)

	// Bind function with partial parameter and return callback with captured partial parameters
	callback2, err := worker.Bind[func(a int) int](add, 10)
	if err != nil {
		panic(err)
	}

	// Run callback with partial parameter
	res2 := callback2.Run(20)
	fmt.Printf("callback2 result: %d\n", res2)

	// Bind function with all parameter and return callback with captured all parameters
	callback3, err := worker.Bind[func() int](add, 100, 200)
	if err != nil {
		panic(err)
	}

	// Run callback without all parameter
	res3 := callback3.Run()
	fmt.Printf("callback3 result: %d\n", res3)
}
