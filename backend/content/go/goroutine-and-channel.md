---
title: "Goroutine 与 Channel 并发模型"
date: "2024-03-01"
tags: ["Go", "并发", "Goroutine", "Channel"]
excerpt: "深入理解 Go 语言的并发哲学：不要通过共享内存来通信，而要通过通信来共享内存。"
author: "camelia"
---

# Goroutine 与 Channel 并发模型

Go 语言的并发模型基于 CSP（Communicating Sequential Processes）理论，核心思想是：

> **Don't communicate by sharing memory; share memory by communicating.**

## Goroutine

Goroutine 是 Go 运行时管理的轻量级线程，创建成本极低（初始栈约 2KB）。

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    go sayHello("camelia") // 启动一个 goroutine
    go sayHello("world")

    time.Sleep(100 * time.Millisecond) // 等待 goroutine 执行
}
```

### WaitGroup 同步

用 `sync.WaitGroup` 比 `time.Sleep` 更可靠：

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Worker %d 开始工作\n", id)
    // 模拟工作...
    fmt.Printf("Worker %d 完成\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait() // 等待所有 worker 完成
    fmt.Println("所有工作完成！")
}
```

## Channel

Channel 是 goroutine 之间通信的管道，类型安全，支持阻塞。

### 无缓冲 Channel

```go
ch := make(chan int) // 无缓冲：发送方阻塞直到接收方准备好

go func() {
    ch <- 42 // 发送（阻塞直到有人接收）
}()

val := <-ch // 接收（阻塞直到有人发送）
fmt.Println(val) // 42
```

### 有缓冲 Channel

```go
ch := make(chan string, 3) // 缓冲大小为 3

ch <- "first"  // 不阻塞
ch <- "second" // 不阻塞
ch <- "third"  // 不阻塞
// ch <- "fourth" // 阻塞！缓冲已满

fmt.Println(<-ch) // "first"
fmt.Println(<-ch) // "second"
```

### Range 遍历 Channel

```go
func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch) // 必须关闭，否则 range 无法结束
}

func main() {
    ch := make(chan int, 5)
    go producer(ch)

    for val := range ch { // 自动检测 channel 关闭
        fmt.Println(val)
    }
}
```

## Select 多路复用

`select` 用于同时等待多个 channel，类似 `switch` 但用于通信：

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() { ch1 <- "来自 ch1" }()
    go func() { ch2 <- "来自 ch2" }()

    for i := 0; i < 2; i++ {
        select {
        case msg := <-ch1:
            fmt.Println("收到:", msg)
        case msg := <-ch2:
            fmt.Println("收到:", msg)
        case <-time.After(1 * time.Second):
            fmt.Println("超时！")
        }
    }
}
```

## 实战：工作池（Worker Pool）

```go
package main

import (
    "fmt"
    "sync"
)

func workerPool(numWorkers, numJobs int) {
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    // 启动 workers
    var wg sync.WaitGroup
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- job * job // 计算平方
            }
        }()
    }

    // 发送任务
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    // 等待所有 worker 完成后关闭结果 channel
    go func() {
        wg.Wait()
        close(results)
    }()

    // 收集结果
    for result := range results {
        fmt.Println(result)
    }
}

func main() {
    workerPool(3, 10) // 3 个 worker，处理 10 个任务
}
```

## 常见陷阱

### Goroutine 泄漏

```go
// ❌ 错误：没有退出机制，goroutine 永远阻塞
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // 没人发送，永远阻塞
        fmt.Println(val)
    }()
}

// ✅ 正确：使用 context 控制生命周期
func noLeak(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done(): // 收到取消信号，退出
            return
        }
    }()
}
```

Go 的并发模型通过 goroutine + channel 让并发编程变得直观而安全喵～
