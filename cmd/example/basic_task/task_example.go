package main

import (
	"context"
	"fmt"
	"time"

	worker "github.com/ISSuh/worker"
)

func taskFunc(index int) {
	fmt.Printf("[taskFunc] index : %d\n", index)
}

func main() {
	// Create task runner with number of worker
	runner := worker.NewTaskRunner(5)

	// run task runner
	// task runner will be stopped when cancel context
	c, cancel := context.WithCancel(context.Background())
	go runner.RunLoop(c)

	for i := 0; i < 50; i++ {
		// Bind task function
		// Task function signature is only can use TaskSigniture. it is func() type
		cb, err := worker.Bind[worker.TaskSigniture](taskFunc, i)
		if err != nil {
			panic(err)
		}

		// Create task
		task := worker.NewTask(cb)

		// Post task to task runner
		runner.PostTask(task)
	}

	time.Sleep(3 * time.Second)

	// when cancel context, task runner will be stopped
	cancel()
}
