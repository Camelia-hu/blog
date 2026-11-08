# enums类型如何表示

### enums类型是什么，有什么用

enums类型是枚举类型，是一组固有的，有名字的常量集合，主要为了实现可读性以及类型安全

### 如何实现

在go里面用const加iota实现

```go
package main

import "fmt"

type Color int

const (
    Red Color = iota
    Green
    Blue
)

func (c Color) String() string {
    return [...]string{"Red", "Green", "Blue"}[c]
}

func main() {
    var c Color = Green
    fmt.Println(c) // 输出: Green
}

```

# 空结构体的作用

### 空结构体大小为0字节，极致省内存

*   作为占位符，结构体里面嵌入一个空结构体字段表示启用某一个功能

```go
type FeatureFlag struct {
    EnableLogging struct{} // 只要这个字段存在，就表示“启用日志”
}

```

*   信号量或者标志

```go
done := make(chan struct{})

// 在某个 goroutine 中
go func() {
    // 完成某些操作
    done <- struct{}{} // 发送信号
}()

// 在主 goroutine 中
<-done // 等待信号

```

*   空结构体channel

```go
// 创建一个缓冲区为 1 的空结构体通道
ch := make(chan struct{}, 1)

// 发送信号
ch <- struct{}{}

// 接收信号
<-ch

```

*   实现集合，相比于bool作为值，更省内存

```go
type Set map[string]struct{}

func main() {
    s := make(Set)
    s["hello"] = struct{}{} // 向集合中添加元素
    fmt.Println(len(s))        // 输出集合的大小
}

```

*   接口的实现

```go
type Printer interface {
    Print()
}

type ConsolePrinter struct{}

func (cp ConsolePrinter) Print() {
    fmt.Println("Hello, world!")
}

```

# go语言字符串转\[]byte会发生内存拷贝吗

会的兄弟，会的

为了实现字符串与byte数组的数据独立性，会为转化的byte数组开辟一个新的内存空间，然后将字符串内容按字节级来复制给byte数组新内存

### 什么是内存拷贝

内存拷贝是一种不同于浅拷贝的，深拷贝是实现内存拷贝的一种方式，内存拷贝颗粒度没有深拷贝那么细，深拷贝是将数据逻辑结构全部复现然后逐个进行内存拷贝，内存拷贝是面对没有那么复杂的数据结构时的深拷贝。

### 字符串的不可变性以及byte数组拷贝的关系

在go语言中字符串类型和byte数组类型有不同的底层实现

字符串类型长度是固定的，内容是不可变的，只可读不可写，在创建后就是一个固定的

```go
type stringStruct struct {
    data *byte  // 指向字符数组
    len  int    // 字符串长度
}

```

byte数组类型长度是可以变化的，内容是可以修改的

```go
type sliceStruct struct {
    data *byte  // 指向底层数组
    len  int    // 数组长度
    cap  int    // 数组容量
}

```

正因为字符串是不可修改的，所以我们为了不修改字符串，在字符串转化为byte数组时就该内存拷贝，以防影响原来的字符串

# 空切片以及nil切片的区别

nil切片没有底层数组，长度和容量也为0

空切片有底层数组，长度和容量为0

```go
var a []int         // nil 切片
b := []int{}        // 空切片
c := make([]int, 0) // 空切片

```

主要要注意json encode的返回结果的区别，空切片是\[]，nil切片是null

# make和new的区别

new可以用在所有数据类型的初始化上，make只能用于map，slice，channel

new返回指针（\*T），make返回引用类型（T）

new只做零值初始化，make做完整结构初始化

```go
m := new(map[string]int)  // 返回 *map，m != nil，但 *m == nil

(*m)["a"] = 1  // ❌ panic：map 结构未初始化

```

```go
m := make(map[string]int)  // 返回 map[string]int

m["a"] = 1  // ✅ 没问题，可直接使用

```

| 比较点          | 零值初始化（`new`）           | 完整结构初始化（`make`）       |
| :----------- | :--------------------- | :-------------------- |
| 内存是否分配       | ✅ 是                    | ✅ 是                   |
| 值内容          | 全为类型的“零值”（0、""、nil）    | 初始化了底层结构，可直接使用        |
| slice 是否可用   | ❌ nil，不能 append        | ✅ 空 slice，可 append    |
| map 是否可用     | ❌ nil，不能插入键值           | ✅ 可立即使用               |
| channel 是否可用 | ❌ nil，不能收发数据           | ✅ 可立即通信               |
| 返回值类型        | `*T`（指针）               | `T` 本身                |
| 推荐用途         | 少用，适用于 struct/new(int) | 常用，适用于 slice/map/chan |

# 可比较类型

*   基本数据类型都是可以比较的
*   指针类型也是可以比较的
*   channel类型也是可以比较的，但是比较的是地址

```go
ch1 := make(chan int)
ch2 := make(chan int)
ch3 := ch1

fmt.Println(ch1 == ch2) // false，不是同一个 channel
fmt.Println(ch1 == ch3) // true，是同一个 channel
//channel的make是开辟一块新内存，所以地址不一样，channel是引用类型，所以浅拷贝，共享底层数据，指向同一块底层数据，所以地址一样
```

*   interface是可以的，只要动态类型和值是可以比较的就行

```go
var i1 interface{} = 10
var i2 interface{} = 10
fmt.Println(i1 == i2) // ✅ true：底层类型和值都一样

var i3 interface{} = []int{1, 2}
var i4 interface{} = []int{1, 2}
// fmt.Println(i3 == i4) // ❌ panic：slice 不可比较

```

*   结构体是可以的，只要所有字段都是可以比较的就行

```go
type Bad struct {
    S []int
}

b1 := Bad{[]int{1, 2}}
b2 := Bad{[]int{1, 2}}
// fmt.Println(b1 == b2) // ❌ 编译错误，因为 slice 不可比较

```

*   数组也是可以的，只要元素类型是可以比较的就行
*   func不行，只能和nil比
*   map不行，只能和nil比，而且，只要可比较类型才能作为map的key
*   slice不行，只能和nil比

