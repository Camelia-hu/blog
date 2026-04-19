# go中如何访问私有成员

私有成员：结构体里面的开头为小写字母的字段

如何访问：只能通过包内访问，公开方法间接访问，以及反射

# 类型断言时候会发生拷贝吗

什么是类型断言：用于将接口类型的值转化为具体类型的值

```markdown
<目标类型的值>，<布尔参数> := <表达式>.( 目标类型 )
```

什么时候发生拷贝：值类型不会发生拷贝，引用类型会

### 详细说一下引用类型和值类型怎么类型断言的

引用类型说白了就是浅拷贝，开辟一块新的内存空间，将原来的引用放进结构体里面，共享底层数据

值类型是开辟一块新的内存空间然后将值赋值给新的内存空间

# go语言接口是怎么实现的

什么是动态类型：接口持有值的实际类型

什么是动态值：接口持有的实际值

接口实际上就是动态类型和动态值的结合

### 接口的底层实现：

空接口：

```go
// 空接口
type eface struct {
    _type *_type         // 数据类型信息
    data  unsafe.Pointer  // 指向实际的数据
}

```

\_type为实际值的类型以及对应信息

data为实际值的指针

非空接口：

```go
// 非空接口
type iface struct {
    tab *itab         // 指向类型信息和方法表
    data unsafe.Pointer // 指向实际的数据
}

```

tab空接口相比多了方法表的指向，因为非空接口对应有方法

data还是一样

\*itab

```go
type itab struct {
    inter *interfacetype // 接口类型信息
    _type *_type         // 实现接口的具体类型信息
    hash  uint32         // 类型 hash 值
    _     [4]byte
    fun   [1]uintptr     // 实现接口方法的函数地址
}

```

inter对应接口本身的信息，如方法名

\_type为接口动态值的实际信息，如int

hash是类型唯一标识，加速map查找

fun的长度是动态变化的，对应接口方法的函数地址

### 接口调用的实际流程：

```go
type Writer interface {
    Write([]byte) (int, error)
}

type MyWriter struct{}

func (mw MyWriter) Write(p []byte) (int, error) {
    fmt.Println("Write called:", string(p))
    return len(p), nil
}

func main() {
    var w Writer = MyWriter{}  // 接口赋值
    w.Write([]byte("hello"))   // 接口调用
}

```

首先生成一个itab：绑定Writer和MyWriter在itab的inter的\_type，将MyWriter.Write方法存进fun里面

然后生成一个w：

```go
iface{
    tab: &itab{...}, // 包含 MyWriter + Writer 的绑定信息
    data: &MyWriter{} // 指向结构体的内存地址
}

```

最后调用w\.Write时：通过w\.tab.fun\[0]找到MyWriter.Write函数，然后将w\.data作为参数进行调用

### 类型断言以及类型判断

类型断言实际上是取出接口的tab信息里面的实际类型信息即\_type进行比对，如果对上了就断言成功

# 怎么实现闭包

闭包是什么：可以引用外部变量的函数值

怎么实现：

```go
package main

import "fmt"

// 实现闭包的函数
func adder() func(int) int {
    sum := 0                   // 定义外部作用域的变量
    return func(x int) int {   // 返回一个匿名函数
        sum += x               // 操作外部作用域的变量
        return sum             // 返回结果
    }
}

func main() {
    pos, neg := adder(), adder()   // 创建两个闭包
    for i := 0; i < 10; i++ {
        fmt.Println(pos(i), neg(-2*i))  // 调用闭包
    }
}

```

# go语言中触发panic的场景有哪些

触发panic会怎样：会非正常中断程序，从下到上依次执行defer，如果捕获到recover就不会中断

1.  数组或者切片越界
2.  空指针解引用
3.  主动抛出panic
4.  类型断言错误
5.  数字错误

    ```go
    result := 1 / 0 // 触发 panic: runtime error: integer divide by zero

    ```
6.  内存越界或者非法操作(一般是使用unsafe包的时候触发)

    ```go
    import "unsafe"

    var p unsafe.Pointer
    *(*int)(p) = 42 // 触发 panic: runtime error: invalid memory address or nil pointer dereference

    ```
7.  运行时错误
8.  使用不安全的库或代码

# go语言中通过指针变量p访问其成员变量title，有哪些方式

一种是使用(\*p).title，另一种是简写p.title
