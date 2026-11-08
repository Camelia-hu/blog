# context作用是什么

一句话总结就是用来传递请求上下文（比如取消信号，超时时间，traceid）等消息的机制，可以在多goroutine之间传递信号

1.  控制并发任务的生命周期，多个goroutine之间共享一个ctx的话，可以实现一处取消，处处取消
2.  传递请求范围内的数据，在处理HTTP请求的时候，context可以用来传递一些请求相关的数据，比如认证，token，用户信息等
3.  设置超时时间，可以通过contexy来设置某个操作的超时时间，一定那超时，相关的goroutine就会取消

## context.Context接口中提供了若干方法

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
```

1.  deadline是用在需要对某个操作设置超时时间的场景
2.  done是实现取消操作的核心方法
3.  err是用于在上层代码中检查和处理不同的取消情况
4.  value，用于传递跨api边界的请求范围数据，例如身份验证令牌，请求标识符等

### 超时控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

resultChan := make(chan string)
go func() {
    time.Sleep(3 * time.Second) // 模拟耗时任务
    resultChan <- "done"
}()

select {
case <-ctx.Done():
    fmt.Println("超时取消:", ctx.Err())
case res := <-resultChan:
    fmt.Println("任务结果:", res)
}
```

### 取消goroutine

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("子任务已取消")
            return
        default:
            fmt.Println("工作中...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}()

time.Sleep(2 * time.Second)
cancel() // 触发取消
time.Sleep(1 * time.Second)
```

### 传递trace*id/user\_id*

```go
func WithTraceID(ctx context.Context, traceID string) context.Context {
    return context.WithValue(ctx, "trace_id", traceID)
}

func GetTraceID(ctx context.Context) string {
    if v := ctx.Value("trace_id"); v != nil {
        return v.(string)
    }
    return ""
}

func Handler(ctx context.Context) {
    fmt.Println("Trace ID:", GetTraceID(ctx))
}

func main() {
    ctx := WithTraceID(context.Background(), "abc-123")
    Handler(ctx)
}

```

# go语言context value的查找过程是怎样的

go会沿着context链从最内层向外逐层查找关键字匹配的值，具体来说

如果当前context的关键词与目标关键词匹配，则返回对应的值

否则，继续向外查找父context，直到匹配到关键字或到达根context，通常是context background

# context如何取消

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println("Context was canceled")
		case <-time.After(5 * time.Second):
			fmt.Println("Context was not canceled within 5 seconds")
		}
	}(ctx)

	cancel()
	time.Slee p(1 * time.Second) // 确保查看结果之前协程有时间执行
}
```

需要注意的是，当一个ctx被cancel后它所有的子ctx都会被cancel

# 什么是反射

略

# 如何利用unsafe获取slice和map的长度

# net/http包中client如何实现长连接

长连接实现的原理是HTTP/1.1通过在HTTP头部设置Connection：keep-alive来告诉服务器希望复用TCP连接以发送和接收后续的HTTP请求与响应，而不是每进行一次HTTP请求就建立一个新连接

在go语言中，要实现长连接可以通过net/http包中的客户端（http.client）

需要关注以下几个参数：

1.  tarnsport参数：这个字段允许你自定义HTTP传输机制，可以配置连接池，代理，tls等，设置最大连接数，代理地址，keepalive等
2.  Timeout参数：这个参数不是直接控制长连接的字段，但它设置了整个请求（抱愧哦连接建立，请求发送，响应接收）的超时时间。合理的超时设置可以避免因为某些原因导致的连接长时间占用

## 怎么关闭连接

如果要主动关闭连接的话，可以这样

```go
    transport, ok := client.Transport.(*http.Transport)
    if ok {
         transport.CloseIdleConnections() // 关闭所有空闲连接。但是
    }

```

### 但是不建议直接关闭连接（net.Conn）

在go的http客户端（http.client）中，你一般拿不到底层的连接对象（net.Conn），除非你自己去实现Transport，而标准库的设计就是为了屏蔽底层连接细节，交给连接池自动管理

go的连接复用是自动管理的，go的http.transport会维护连接池（Connection Pool），它通过复用连接来提升性能。

1.  首次请求时：创建新连接（net.Conn），放入连接池中（按照目标主机地址分）
2.  后续请求：检查连接池中是否已有空闲连接->复用它（避免TCP/TLS握手），如果没有空闲连接就再新建连接
3.  请求完成，如果设置了Keep-Alive，连接不会关闭而是放回池子
4.  主动关闭连接池：transport.CloseIdConnections（）

什么是连接复用

连接复用是指多个http请求公用同一个TCP连接（在Keep-Alive模式上）

避免每次都三次握手（TCP），握手（TLS）。减少延迟，提升吞吐量。节省系统资源

# net/http如何做连接池

在net/http包中，连接池是通过http.Transport来实现的。http.Transport默认提供了一个内置的连接池，它会为每个主机保存多个空闲连接，并在发送请求时复用这些连接。这样可以避免每次请求都建立新的TCP连接，提升性能

## 连接池的工作原理

*   http.Transport会为每个主机维护一个连接池
*   当发送请求时，Transport会产生复用已有的连接，如果有空闲连接则复用它
*   如果没有空闲连接，Transport会建立新的连接，并将其放入连接池中

## 核心参数

*   MaxIdleConns：设置最大空闲数，控制连接池的大小
*   MaxIdleConnsPerHost：设置每个主机的最大空闲连接数
*   IdleConnTimeout：设置空闲连接的最大保持时间，如果连接在指定时间内没有使用，它会被关闭

## http.Transport和连接池的实现实例

http.Transport负责管理与远程服务器的连接。它会在内部为每个主机维护一个连接池，避免频繁的建立与关闭连接

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func main() {
    // 配置自定义 Transport 以启用连接池
    transport := &http.Transport{
        MaxIdleConns:        100,                          // 最大空闲连接数
        MaxIdleConnsPerHost: 10,                           // 每个主机的最大空闲连接数
        IdleConnTimeout:     30 * time.Second,             // 空闲连接最大保持时间
    }

    client := &http.Client{
        Transport: transport,
    }

    // 发送请求
    resp, err := client.Get("http://example.com")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Response status:", resp.Status)
}

```

连接池是按照主机划分的，即协议+ip+端口只要有一个不同就是不同主机，就会建立不同连接池

## 部分实现原理的源码解析

http.Transport的源代码主要作用是

*   处理HTTP请求的发送
*   管理连接池
*   支持持久连接
*   支持并发请求复用连接

并发请求复用连接

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"
)

func main() {
    client := &http.Client{
        Transport: &http.Transport{
            MaxIdleConns:        10,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     30 * time.Second,
        },
        Timeout: 10 * time.Second,
    }

    var wg sync.WaitGroup
    url := "https://httpbin.org/get"

    concurrency := 5
    wg.Add(concurrency)

    start := time.Now()

    for i := 0; i < concurrency; i++ {
        go func(id int) {
            defer wg.Done()

            resp, err := client.Get(url)
            if err != nil {
                fmt.Printf("Goroutine %d: request error: %v\n", id, err)
                return
            }
            defer resp.Body.Close()

            // 读取响应体，必须读完或关闭，否则连接无法复用
            _, _ = io.Copy(io.Discard, resp.Body)

            fmt.Printf("Goroutine %d: status=%s\n", id, resp.Status)
        }(i + 1)
    }

    wg.Wait()
    fmt.Printf("All requests done in %v\n", time.Since(start))
}

```

你会发现这些请求实际上只建立了一两条tcp连接，后续请求直接复用

### 连接池的结构

http.Tranport内部用coonPool（连接池）来管理所有连接，为了高效存储和查找连接，go使用了sync.Map来存储不同目标主机的空闲连接

```go
// Transport is an HTTP handler that implements RoundTripper.
type Transport struct {
    // Base settings for Transport.
    // Transport implements RoundTripper
    proxy             func(*http.Request) (*url.URL, error)
    dial              func(ctx context.Context, network, addr string) (net.Conn, error)
    dialTLS           func(ctx context.Context, network, addr string) (net.Conn, error)
    tlsClientConfig   *tls.Config
    disableCompression bool

    // The connPool is a map of host -> connection pool.
    connPool          *connPool
}

```

connPool是一个自定义的结构体，用于存储和管理不同主机之间的连接池

### connPool的结构与管理

connPool负责管理空闲连接的复用与清理，保存了每个主机的连接池

```go
type connPool struct {
    mu sync.Mutex // protects the following fields
    m  map[string]*httpConnPool
}

type httpConnPool struct {
    mu        sync.Mutex
    maxIdle   int // 最大空闲连接数
    idleConns []*httpConn // 存储空闲连接的列表
}

type httpConn struct {
    conn    net.Conn
    lastUsed time.Time // 记录上次使用时间
}

```

connPool：维护了一个哈希表，键为目标主机（通常为主机名或ip地址），值为该主机对应的连接池（httpConnPool）

httpConnPool：每个主机对应的连接池，它维护了一组空闲连接idleConns，以及最大空闲连接数maxIdle

httpConn：每个连接的结构体，包含一个net.Conn和lastUsed字段来记录连接的最后使用时间

### 连接复用过程

主要在getConn和putConn方法中

#### getConn，获取连接

```go
func (p *connPool) getConn(network, addr string, ctx context.Context) (conn net.Conn, err error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    // 获取连接池
    pool, exists := p.m[addr]
    if !exists {
        pool = &httpConnPool{maxIdle: defaultMaxIdleConns}
        p.m[addr] = pool
    }

    // 从连接池中取出一个空闲连接
    if len(pool.idleConns) > 0 {
        conn = pool.idleConns[len(pool.idleConns)-1]
        pool.idleConns = pool.idleConns[:len(pool.idleConns)-1]
        return conn, nil
    }

    // 如果没有空闲连接，创建一个新的连接
    conn, err = p.dial(ctx, network, addr)
    if err != nil {
        return nil, err
    }

    return conn, nil
}
```

1.  首先检查是否有该主机的连接池，如果没有就新建一个
2.  再去看连接池里面有没有空闲连接，有就取出来，没有就新建

#### putConn，归还连接

```go
func (p *connPool) putConn(addr string, conn net.Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()

    // 获取目标主机的连接池
    pool, exists := p.m[addr]
    if !exists {
        pool = &httpConnPool{maxIdle: defaultMaxIdleConns}
        p.m[addr] = pool
    }

    // 将连接放入连接池中
    if len(pool.idleConns) < pool.maxIdle {
        pool.idleConns = append(pool.idleConns, &httpConn{conn: conn, lastUsed: time.Now()})
    } else {
        conn.Close() // 如果连接池已满，则关闭该连接
    }
}

```

将连接归还到目标主机的连接池里面，如果达到最大空闲连接数就关闭连接，没有就塞进去
