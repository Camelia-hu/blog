---
title: "Interface 与泛型：Go 的抽象之道"
date: "2024-03-15"
tags: ["Go", "Interface", "泛型", "设计模式"]
excerpt: "理解 Go 的接口隐式实现机制与 1.18 引入的泛型，写出更灵活、可复用的代码。"
author: "camelia"
---

# Interface 与泛型：Go 的抽象之道

## Interface（接口）

Go 的接口是**隐式实现**的——无需显式声明 `implements`，只要实现了接口的所有方法，就自动满足该接口。

```go
type Animal interface {
    Sound() string
    Name() string
}

type Dog struct{}
type Cat struct{}

func (d Dog) Sound() string { return "汪汪" }
func (d Dog) Name() string  { return "小狗" }

func (c Cat) Sound() string { return "喵喵" }
func (c Cat) Name() string  { return "小猫" }

func greet(a Animal) {
    fmt.Printf("%s 说：%s\n", a.Name(), a.Sound())
}

func main() {
    greet(Dog{}) // 小狗 说：汪汪
    greet(Cat{}) // 小猫 说：喵喵
}
```

## 空接口与类型断言

```go
// interface{} 或 any 可接受任意类型
func printAny(v any) {
    switch val := v.(type) {
    case int:
        fmt.Printf("整数: %d\n", val)
    case string:
        fmt.Printf("字符串: %s\n", val)
    case []int:
        fmt.Printf("整数切片: %v\n", val)
    default:
        fmt.Printf("未知类型: %T\n", val)
    }
}
```

## 接口组合

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}
```

## 泛型（Go 1.18+）

泛型让函数和数据结构能处理多种类型，同时保持类型安全。

### 泛型函数

```go
// 没有泛型时，需要为每种类型写一个函数
func MaxInt(a, b int) int {
    if a > b { return a }
    return b
}

// 有泛型后，一个函数搞定
func Max[T int | float64 | string](a, b T) T {
    if a > b { return a }
    return b
}

func main() {
    fmt.Println(Max(3, 5))       // 5
    fmt.Println(Max(3.14, 2.71)) // 3.14
    fmt.Println(Max("go", "java")) // go（字典序）
}
```

### 类型约束

```go
import "golang.org/x/exp/constraints"

// constraints.Ordered 包含所有可比较排序的类型
func Min[T constraints.Ordered](a, b T) T {
    if a < b { return a }
    return b
}

// 自定义约束
type Number interface {
    int | int8 | int16 | int32 | int64 |
        float32 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    fmt.Println(Sum([]int{1, 2, 3, 4, 5}))       // 15
    fmt.Println(Sum([]float64{1.1, 2.2, 3.3}))   // 6.6
}
```

### 泛型数据结构：栈

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    last := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return last, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

func main() {
    s := &Stack[string]{}
    s.Push("hello")
    s.Push("world")

    if val, ok := s.Pop(); ok {
        fmt.Println(val) // "world"
    }
}
```

## 最佳实践

- **接口尽量小**：`io.Reader`（1个方法）比大接口更灵活
- **面向接口编程**：函数参数用接口而非具体类型
- **泛型适度使用**：简单场景用 `any` + 类型断言即可，复杂复用场景再用泛型
- **接口定义在使用方**：不要在实现包中定义接口
