# 文件资源不关闭会发生什么事

1.  操作系统会给每个打开的文件资源分配一个文件描述符（FD），每个操作系统能打开的FD是有限的，为1024个，如果不关闭的话会产生错误
2.  如果是写操作，不关闭的话缓冲区的对应修改可能无法刷新到对应的磁盘
3.  操作系统不止会维护一个FD，还会维护一个内核态资源，如文件表项，如果不释放就会占用资源
4.  如果文件加锁，但是不关闭，就会造成其他程序无法访问，造成死锁

### 打开文件的流程

用户调用os.open，系统调用syscall.open，内核分配文件描述符创建file table记录加载inota信息

# go语言的局部变量是分配在栈上还是堆上

都可以，主要看是否逃离函数作用域

### 栈和堆的区别

栈：

栈是线程私有的，分配和释放内存由编译器管理

分配速度快，通常用于局部变量的分配

堆：

堆是全局共享的，由GC管理

堆适用于生命周期长且动态分配的内存，但速度较慢

栈和堆两者的结合设计主要是为了保证灵活性以及速度的合适调度

### 逃逸分析的工作原理

##### 逃逸分析是什么

逃逸分析是一种编译时优化技术，编译器通过分析变量的使用场景来决定变量的生命周期以及内存是分配在栈上还是堆上，核心逻辑是判断变量的引用是否超出函数的生命周期

触发堆分配的场景

返回指针：

```go
func createPointer() *int {
    x := 10
    return &x // x 逃逸到堆上
}

```

闭包捕获：

```go
func closure() func() int {
    x := 10
    return func() int { // x 逃逸到堆上
        return x
    }
}

```

跨goroutine传递：

```go
func goroutine() {
    x := 10
    go func() {
        fmt.Println(x) // x 逃逸到堆上
    }()
}

```

# 数组和切片的区别

数组是值传递，作为参数传递时不会影响到原数组，切片虽然也是值传递，其构成为指向底层数组的指针，容量以及当前长度，传递的值是指针，所以有类引用效果

当切片扩容后会开辟一个新的底层数组，然后将原来底层数组复制过来，append是返回一个新切片，如果扩容则底层数组是新的。在切片长度小于256时容量扩张为两倍，超过256后增长四分之三

# 切片作为函数参数时要注意的点

1.  切片作为参数时for循环只能用s\[i]修改，而不能用for \_, v := range s来修改，因为v是copy过来的
2.  如果在函数内部对切片进行增加以及删除时，可以将切片指针传进来，不然外面的切片还是原来的切片
3.  切片是引用类型，传进来的是对底层数组的指针copy
4.  由于切片具有动态大小，所以你增加或者删除元素时可能会引发指向的底层数组地址重新分配
5.  最好通过返回新切片，如果你要增加或者减少切片元素时
6.  切片可以用const或者只读

# go中的rune类型是什么

’rune类型实质上是int32的别名，主要作用是可以表示中文或其他字符，如“你好hello”用rune类型来循环读取就是你，好，h，l，l，o，用\[]byte类型读就是乱码。for range读取的是rune类型，为了方便中文或其他字符的处理

# 深拷贝和浅拷贝

浅拷贝是指将值拷贝过来，包括内部结构们内部结构还是原来的值，修改引用会影响浅拷贝的值

深拷贝是指将值拷贝在新的内部结构上，内存结构是新分配的，修改引用不会影响深拷贝的值

浅拷贝的一般是较为复杂的类型，如     切片，channel，interface，map，指针，func

深拷贝的一般是基础数据类型，如int，string等

```go
package main

import "fmt"

func main() {
    // 原始二维切片
    original := [][]int{
        {1, 2},
        {3, 4},
    }

    // -------------------
    // 浅拷贝
    // -------------------
    shallow := original

    // 改变浅拷贝的数据
    shallow[0][0] = 100

    fmt.Println("原始 original:", original)
    fmt.Println("浅拷贝 shallow:", shallow)

    // -------------------
    // 深拷贝
    // -------------------
    deep := make([][]int, len(original))
    for i := range original {
        deep[i] = make([]int, len(original[i]))
        copy(deep[i], original[i])
    }

    // 改变深拷贝的数据
    deep[0][0] = 999

    fmt.Println("原始 original:", original)
    fmt.Println("深拷贝 deep   :", deep)
}

```

# go中的map是不可寻址的，那怎么修改值

如果值是基本类型的话，那么可以直接修改：

```go
m := make(map[string]string)
m["key1"] = "mianshi"
m["key1"] = "mianshiya" //直接修改

```

如果值是指向结构体的指针类型也可以直接修改：

```go
type Person struct {
   Name string
   Age  int
}

m := make(map[int]*Person)
m[1] = &Person{"Alice", 30}
m[1].Age = 31  // 修改结构体的属性
fmt.Println(m[1].Age)  // 输出 31

```

如果值是结构体类型的话要先创建一个结构体temp，然后替代原来的结构体值来修改：

```go
m := make(map[int]Person)
m[1] = Person{"Alice", 30}

// 修改中间变量
temp := m[1]
temp.Age = 31
m[1] = temp  // 将修改后的值放回 map
fmt.Println(m[1].Age)  // 输出 31


```

### 为什么map的值不可寻址

因为map是线程不安全的，怕发生竞态

```go
m := make(map[int]int)
m[1] = 10
p := &m[1]  // 错误：无法取map值的地址

```

# 有类型常量和无类型常量

常量是程序运行中无法改变的量，就是const

有类型常量就是你在const的时候定义了它的类型

```go
const pi float64 = 3.14  // 有类型常量

```

无类型常量就是你没有定义

```go
const pi = 3.14  // 无类型常量

```

无类型常量由编译器自动推断类型，有类型常量无法修改。前者保证灵活性，后者保证类型安全

# 函数传参时倾向于传切片而非数组，为什么

1.  效率更高，因为切片传的是底层数组的引用，数组是传来整个数组的copy，数组开销太大了
2.  灵活性更强，可以动态改变长度
3.  内置功能更丰富，silce的生态更丰富

# go中的引用类型和指针

引用类型实际上就是一个包括指针以及其他元素的结构体如切片，指针就是指针

```go
type Slice struct {
    array unsafe.Pointer // 指向底层数组的指针
    len   int            // 切片的长度
    cap   int            // 切片的容量
}

```

引用类型通过值传递来作为参数，而指针类型一般是为了避免拷贝整个对象带来的额外开销而使用

```go
type Person struct {
    Name string
    Age  int
}

func updatePerson(p *Person) {
    p.Name = "Alice" // 修改的是指针指向的数据
}

func main() {
    p := &Person{Name: "Bob", Age: 30}
    updatePerson(p)
    fmt.Println(p.Name) // 输出 Alice
}

```

# go语言map如何访问

一种方式是

```go
value := myMap[key]

```

这样的话如果不存在就会返回零值，如int会返回0

另一种方式是

```go
value, exists := myMap[key]

```

这样通过exist来判断是否存在

如果是nil map（即没有初始化分配内存）的话，直接访问不会产生panic，而会返回零值以及false

```go
var nilMap map[string]int
value, ok := nilMap["key"]
fmt.Println(value, ok)  // 输出: 0 false

```

# map的无序性

```go
myMap := map[string]int{
    "apple":  3,
    "banana": 2,
    "cherry": 5,
}

for key, value := range myMap {
    fmt.Println(key, value)
}

```

这样每次打印出来的顺序可能是不一样的

# defer的执行顺序是怎样的

后进先出

### 使用defer的注意点

1.  defer不能跨goroutine捕捉panic
2.  defer函数下的参数会在defer执行时被清算，而不是最后执行时

```go
func main() {
    a := 10
    defer fmt.Println(a) // 立刻捕获了 a 的值
    a = 20
}

```

这里输出10而不是20

# tag的用法

tag是结构体字段后面加的元数据字符串，描述结构体字段的附加信息，通常用于与外部系统交流沟通

常见用途：

1.  json编码解码
2.  数据库orm
3.  表单验证
4.  xml编码解码

# 如何比较两个切片是否一样

自定义循环或者用reflect.DeepEqual，这个函数可以比较复杂数据结构是否一样（切片，映射，结构体等）

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    slice1 := []string{"apple", "banana", "cherry"}
    slice2 := []string{"apple", "banana", "cherry"}
    slice3 := []string{"apple", "banana"}

    fmt.Println(reflect.DeepEqual(slice1, slice2)) // true
    fmt.Println(reflect.DeepEqual(slice1, slice3)) // false
}

```

# go中常见占位符

%v与%+v的区别在于打印结构体时%+v不仅会打印值也会打印字段

```go
package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	p := Person{"Alice", 30}
	fmt.Printf("%v\n", p)  // 输出：{Alice 30}
	fmt.Printf("%＋ v\n", p) // 输出：{Name:Alice Age:30}
}

```

%T打印变量类型

```go
fmt.Printf("%T\n", p)  // 输出：main.Person

```

%q打印带引号的字符串

```go
fmt.Printf("%q\n", "hello") // 输出："hello"

```

其他

```go
var num = 42
fmt.Printf("%d\n", num) // 输出：42
fmt.Printf("%s\n", "hello") // 输出：hello
fmt.Printf("%f\n", 3.14) // 输出：3.140000

```

