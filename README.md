# worker

**worker** is a Go package that implements a task-based worker system inspired by the internal task-based operations of the **Chromium** browser. This package provides functionality to add tasks to a queue and process them with multiple worker threads.

## Features

- Callback Binding
- Task Runner

## Usage

```sh
go get github.com/ISSuh/worker
```

### Callback & Bind

- [callback_example.go](./cmd/example/callback/callback_example.go)

```go
package main

import (
    "context"
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
```
```bash
callback1 result: 3
callback2 result: 30
callback3 result: 300
```

### basic Task

- [task_example.go](./cmd/example/basic_task/task_example.go)

```go
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
```
```bash
[taskFunc] index : 2
[taskFunc] index : 4
[taskFunc] index : 5
[taskFunc] index : 6
...
```

### delay Task

- [delay_task_example.go](./cmd/example/delay_task/delay_task_example.go)

```go
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

```
```bash
[taskFunc] index : 0, duration : 0
[taskFunc] index : 10, duration : 10000000
[taskFunc] index : 20, duration : 20000000
[taskFunc] index : 30, duration : 30000000
[taskFunc] index : 40, duration : 40000000
```