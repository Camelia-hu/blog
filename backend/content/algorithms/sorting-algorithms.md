---
title: "排序算法详解"
date: "2024-01-15"
tags: ["算法", "排序", "时间复杂度"]
excerpt: "深入探讨各种排序算法的原理、实现与性能分析，包括快速排序、归并排序、堆排序等经典算法。"
author: "Mayo"
---

# 排序算法详解

排序算法是计算机科学中最基础且重要的算法之一。本文将深入探讨几种经典排序算法的原理与实现。

## 时间复杂度对比

| 算法 | 最好 | 平均 | 最坏 | 空间 | 稳定 |
|------|------|------|------|------|------|
| 快速排序 | O(n log n) | O(n log n) | O(n²) | O(log n) | 不稳定 |
| 归并排序 | O(n log n) | O(n log n) | O(n log n) | O(n) | 稳定 |
| 堆排序 | O(n log n) | O(n log n) | O(n log n) | O(1) | 不稳定 |
| 插入排序 | O(n) | O(n²) | O(n²) | O(1) | 稳定 |
| 冒泡排序 | O(n) | O(n²) | O(n²) | O(1) | 稳定 |

## 快速排序

快速排序采用**分治**策略，通过选取基准元素将数组分为两部分递归排序。

```go
func quickSort(arr []int, low, high int) {
    if low < high {
        pivot := partition(arr, low, high)
        quickSort(arr, low, pivot-1)
        quickSort(arr, pivot+1, high)
    }
}

func partition(arr []int, low, high int) int {
    pivot := arr[high]
    i := low - 1
    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}
```

**特点：**
- 平均时间复杂度 O(n log n)，实际表现非常好
- 最坏情况（已排序数组）退化为 O(n²)
- 原地排序，空间复杂度 O(log n)（递归栈）

## 归并排序

归并排序采用分治策略，将数组分成两半分别排序后再合并。

```go
func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])
    right := mergeSort(arr[mid:])
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}
```

**特点：**
- 时间复杂度稳定在 O(n log n)
- 需要额外 O(n) 空间
- 稳定排序，适合外部排序

## 堆排序

利用堆数据结构进行排序，分为建堆和排序两个阶段。

```go
func heapSort(arr []int) {
    n := len(arr)
    // 建大根堆
    for i := n/2 - 1; i >= 0; i-- {
        heapify(arr, n, i)
    }
    // 逐步提取最大值
    for i := n - 1; i > 0; i-- {
        arr[0], arr[i] = arr[i], arr[0]
        heapify(arr, i, 0)
    }
}

func heapify(arr []int, n, i int) {
    largest := i
    left, right := 2*i+1, 2*i+2
    if left < n && arr[left] > arr[largest] {
        largest = left
    }
    if right < n && arr[right] > arr[largest] {
        largest = right
    }
    if largest != i {
        arr[i], arr[largest] = arr[largest], arr[i]
        heapify(arr, n, largest)
    }
}
```

## 总结

- **小规模数据**：插入排序效率更高，常数因子小
- **大规模数据**：快速排序平均性能最佳
- **需要稳定性**：选择归并排序
- **内存受限**：堆排序是最佳选择，O(1) 额外空间
- **外部排序**：归并排序天然适合
