---
title: "进程与线程管理"
date: "2024-02-10"
tags: ["操作系统", "进程", "线程", "并发", "同步"]
excerpt: "深入理解操作系统中进程和线程的管理机制，包括调度算法、同步原语和死锁处理。"
author: "Mayo"
---

# 进程与线程管理

操作系统的核心任务之一是管理进程和线程，为程序提供执行环境。

## 进程 vs 线程

| 特性 | 进程 | 线程 |
|------|------|------|
| 地址空间 | 独立（隔离） | 共享进程地址空间 |
| 创建开销 | 大（需复制地址空间） | 小 |
| 切换开销 | 大（需切换地址空间） | 小 |
| 通信方式 | IPC（管道、信号、共享内存） | 直接共享内存 |
| 安全性 | 高（进程间隔离） | 低（线程间无隔离）|
| 崩溃影响 | 不影响其他进程 | 可能导致整个进程崩溃 |

## 进程状态机

```
新建 ──→ 就绪 ──→ 运行
              ↑      │
              │      ↓
           等待 ←── 阻塞

              └──→ 终止
```

- **新建**：进程正在被创建
- **就绪**：等待 CPU 调度
- **运行**：正在 CPU 上执行
- **阻塞**：等待 I/O 或事件
- **终止**：执行完毕

## 进程调度算法

### FCFS（先来先服务）

- 按到达顺序分配 CPU
- 简单公平，但可能导致**护航效应**（Convoy Effect）
- 长进程阻塞短进程

### SJF（短作业优先）

- 优先执行预计执行时间最短的进程
- 最优平均等待时间
- **缺点**：需要预知执行时间，可能导致长进程饥饿

### 轮转调度（Round Robin）

- 每个进程执行固定**时间片**（通常 10-100ms）
- 超时则切换到下一个就绪进程
- 适合交互式系统，响应时间可预测

### 优先级调度

- 为进程分配优先级，高优先级先执行
- 可能导致低优先级进程**饥饿**
- 解决方案：**老化**（Aging）—— 随等待时间提升优先级

## 同步原语

### 互斥锁（Mutex）

```go
var mu sync.Mutex

func criticalSection() {
    mu.Lock()
    defer mu.Unlock()
    // 临界区：同时只有一个 goroutine 执行
    sharedResource++
}
```

### 读写锁（RWMutex）

适合**读多写少**的场景：

```go
var rwmu sync.RWMutex

// 读操作（允许并发）
func read() int {
    rwmu.RLock()
    defer rwmu.RUnlock()
    return sharedData
}

// 写操作（独占）
func write(val int) {
    rwmu.Lock()
    defer rwmu.Unlock()
    sharedData = val
}
```

### 信号量（Semaphore）

```go
// 使用 channel 模拟信号量
sem := make(chan struct{}, 5) // 最多 5 个并发

func limitedConcurrency() {
    sem <- struct{}{} // 获取许可
    defer func() { <-sem }() // 释放许可
    // 执行任务
}
```

### 条件变量（Cond）

```go
var mu sync.Mutex
cond := sync.NewCond(&mu)
queue := []int{}

// 生产者
func producer() {
    mu.Lock()
    queue = append(queue, 1)
    cond.Signal() // 通知消费者
    mu.Unlock()
}

// 消费者
func consumer() {
    mu.Lock()
    for len(queue) == 0 {
        cond.Wait() // 等待通知，自动释放锁
    }
    item := queue[0]
    queue = queue[1:]
    mu.Unlock()
    _ = item
}
```

## 死锁

### 死锁的四个必要条件

1. **互斥**：资源同时只能被一个进程使用
2. **占有并等待**：进程占有资源同时等待其他资源
3. **不可抢占**：资源不能被强制夺走
4. **循环等待**：存在进程的等待环路

### 死锁预防

破坏任意一个必要条件即可：

- **破坏互斥**：使用无锁数据结构
- **破坏占有并等待**：一次性申请所有资源
- **破坏不可抢占**：允许高优先级抢占低优先级资源
- **破坏循环等待**：对资源编号，按序申请（**最常用**）

```go
// 按固定顺序加锁，避免死锁
func transfer(from, to *Account, amount int) {
    // 始终先锁 ID 小的账户
    if from.ID < to.ID {
        from.mu.Lock()
        to.mu.Lock()
    } else {
        to.mu.Lock()
        from.mu.Lock()
    }
    defer from.mu.Unlock()
    defer to.mu.Unlock()
    // 执行转账
}
```

## 上下文切换

上下文切换是操作系统切换进程/线程时保存和恢复状态的过程，包括：

- CPU 寄存器（PC、SP、通用寄存器）
- 进程控制块（PCB）
- 页表（进程切换时）
- TLB 刷新（进程切换时，开销最大）

**优化建议**：减少锁竞争、使用协程（goroutine）代替线程、使用无锁算法。
