# map

## 底层实现

```go
type hmap struct {
	count int // map中元素的个数
	flags uint8 // 状态标志位，标记map的一些状态
	B uint8 // 桶数以2为底的对数
	noverflow uint16 // 溢出桶的数量近似值
	hash0 uint32 // 哈希种子

	bucket unsafe.Pointer // 指向bucket数组的指针
	oldbuckets unsafe.Pointer // 是一个指向buckets数组的指针，在扩容时，指向老的bucket数组（大小为新数组的一半），非扩容时，oldbuckets为空
	nevacuate uintptr // 表示扩容进度的一个计数器，小于该值的桶已经完成迁移
	
	extra *mapextra // 指向mapextra 的指针，mapextra 存储map中的溢出桶
}
```

## 核心字段详解

### count

元素数量

比如：

```go
m := map[string]int{
"a":1,
"b":2,
}
```

count=2

len(map) 直接读这里。

时间复杂度：

O(1)



### B

决定桶数量。

公式：

bucket 数 = 2^B

比如：

B=0 → 1 bucket
B=1 → 2 bucket
B=2 → 4 bucket
B=10 → 1024 bucket

### buckets

指向当前桶数组。

比如：

bucket0
bucket1
bucket2
...
bucketN

都在连续内存。



### oldbuckets

扩容时旧桶数组。

因为 Go map 采用：

渐进迁移

所以新旧桶会同时存在一段时间。



### nevacuate

记录迁移进度。

比如旧桶有 1024 个：

已经迁了 300 个：

nevacuate=300

### hash0

随机 hash seed。

用于 hash 扰动。

目的是：

防 hash 冲突攻击。

比如用户构造：

100万个key
全部落一个桶

复杂度退化：

O(n)

seed 后攻击难度大很多。



## bucket结构

    tophash[8]

    key[8]

    value[8]

    overflow ptr

    +----------------+
    | tophash[8]     |
    +----------------+
    | key0           |
    | key1           |
    | ...            |
    | key7           |
    +----------------+
    | value0         |
    | value1         |
    | ...            |
    | value7         |
    +----------------+
    | overflow ptr   |
    +----------------+



一个 bucket 放：8组 kv

这是固定设计。

为什么是 8？

runtime 做过平衡：

太小：

桶太多

cache miss 多

太大：

桶太胖

扫描慢

8 是工程平衡点。



## tophash是什么

不是完整 hash。只是高 8 bit。

比如：

完整 hash：01101011 10010101 ...

高位：01101011

存入：tophash

作用：快速过滤。

查找时：

先比 top hash

相同再比 key

减少大量比较。

像布隆过滤的第一层筛选。



## 查找流程

第一步：算 hash

调用类型自己的 hash 算法：比如 string：扫描字符串字节。得到：hash(key)
第二步：定位 bucket

公式：

bucketIndex = hash & (2^B-1)

比如：16 个桶：mask=15。hash：101011011。低4位：1011。就是：11号桶。因为：& 比 % 快&#x20;

第三步：匹配 tophash

遍历：tophash\[0:8]。快速比。不匹配直接跳过。

第四步：比较 key

如果 top 相同：再比较 key。string 逐字节比较。int 直接比较。interface 要走动态比较器。

第五步：取 value

返回对应位置 value。

复杂度：

平均：O(1)

冲突严重：O(n)



## 哈希冲突

两个 key 落同桶。

比如：bucket3。已经 8 个元素。第 9 个怎么办？

挂 overflow bucket。

形成链：

bucket3
↓
overflow1
↓
overflow2

继续放。

这叫：溢出桶链

问题：

链长后：

查找退化。

CPU cache miss 激增。

性能下降明显。



## 插入流程 m\[k]=v

步骤：

算 hash。定位桶。找空位。扫描 8 格。找到 empty slot。插入。满了，申请 overflow bucket。挂链。插进去。count++，数量+1。判断扩容，如果触发就 grow。



## 删除 delete(m,k)

删除不会真的搬迁元素。只做标记。

类似：occupied → emptyOne

表示：这里删过。

为什么？

因为开放寻址/扫描逻辑依赖状态连续性。

直接清空会破坏探测。

value 清零：帮助 GC。count--



## 扩容机制（核心）

Go map 最精华部分。不是一次性 rehash。是：渐进扩容

触发有两种情况：

### 负载因子过高

load factor：count / bucket数，超过约：6.5就扩容。

比如：1024 bucket放：7000+就触发。

### overflow 太多

虽然元素不多。冲突严重。overflow 很长。也扩容。

叫：

sameSizeGrow

同大小整理。

重新洗牌。

减少冲突链。



## 渐进迁移

不是一次搬完。因为 pause 太大。一次搬 1百万元素：延迟爆炸。

Go 做法：边用边搬。

新建：2倍 buckets。旧桶挂：oldbuckets

每次：插入， 删除， 查找的时候顺便迁一点。

最终：oldbuckets=nil。结束。这就是：incremental rehash



## map遍历无序是因为bucket遍历起点随机

## map不能取地址是因为扩容的时候地址会变

## map并发panic是因为扩容迁移的时候bucket状态变化容易读到脏数据，flag记录状态

## sync.Map存在是因为适配读多写少的场景，普通map+RWMutex适合写多

## 大map对gc压力很大





































&#x20;
