package main

import (
	"context"
	"fmt"
	"time"

	worker "github.com/ISSuh/worker"
)

func taskFunc(index int, duration time.Duration) {
	fmt.Printf("[taskFunc] index : %d, duration : %d\n", index, duration)
}

func main() {
	// Create task runner with number of worker
	runner := worker.NewTaskRunner(5)

	// run task runner
	// task runner will be stopped when cancel context
	c, cancel := context.WithCancel(context.Background())
	go runner.RunLoop(c)

	for i := 0; i < 50; i += 10 {
		duration := time.Duration(i) * time.Millisecond

		// Bind task function
		// Task function signature is only can use TaskSigniture. it is func() type
		cb, err := worker.Bind[worker.TaskSigniture](taskFunc, i, duration)
		if err != nil {
			panic(err)
		}

		// Create delay task with duration
		delayTask := worker.NewDelayTask(duration, cb)

		// Post task to task runner
		runner.PostTask(delayTask)
	}

	time.Sleep(3 * time.Second)

	// when cancel context, task runner will be stopped
	cancel()
}
