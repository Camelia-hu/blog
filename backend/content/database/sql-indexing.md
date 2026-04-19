---
title: "SQL 索引深度解析"
date: "2024-01-20"
tags: ["数据库", "MySQL", "索引", "性能优化"]
excerpt: "深入理解 B+ 树索引的原理，掌握索引优化技巧，提升数据库查询性能。"
author: "Mayo"
---

# SQL 索引深度解析

索引是数据库性能优化中最重要的工具之一。理解索引原理是每个后端开发者的必修课。

## B+ 树结构

MySQL InnoDB 使用 B+ 树作为索引结构，具有以下特点：

- 所有**数据存储在叶节点**，非叶节点只存储键值
- 叶节点通过**双向链表**相连，支持高效范围查询
- 树高通常为 3-4 层，每次查询 I/O 次数极少
- 非叶节点存放更多键值，降低树高

## 索引类型

### 聚集索引（Clustered Index）

主键索引就是聚集索引，数据按主键顺序物理存储在磁盘上。

```sql
CREATE TABLE users (
    id INT PRIMARY KEY,  -- 聚集索引
    name VARCHAR(100),
    email VARCHAR(200)
);
```

> **注意**：每张表只能有一个聚集索引，数据文件就是索引文件。

### 二级索引（Secondary Index）

```sql
CREATE INDEX idx_email ON users(email);
```

二级索引的叶节点存储的是**主键值**，查询时需要**回表**操作：

1. 通过二级索引找到主键值
2. 再通过主键在聚集索引中查找完整行数据

### 联合索引（Composite Index）

```sql
CREATE INDEX idx_name_age ON users(name, age);
```

## 索引优化技巧

### 最左前缀原则

对于联合索引 `(a, b, c)`，查询必须从最左列开始：

```sql
-- ✅ 可以使用索引
SELECT * FROM t WHERE a = 1;
SELECT * FROM t WHERE a = 1 AND b = 2;
SELECT * FROM t WHERE a = 1 AND b = 2 AND c = 3;
SELECT * FROM t WHERE a = 1 AND c = 3;  -- 只用到 a 部分

-- ❌ 不能使用索引
SELECT * FROM t WHERE b = 2;         -- 跳过了 a
SELECT * FROM t WHERE b = 2 AND c = 3;  -- 跳过了 a
```

### 覆盖索引

查询的列全部包含在索引中，**无需回表**，性能极佳：

```sql
-- 假设有索引 (name, age)
-- 覆盖索引查询，Extra 显示 Using index
SELECT name, age FROM users WHERE name = 'Alice';

-- 需要回表（SELECT *）
SELECT * FROM users WHERE name = 'Alice';
```

### 索引失效场景

```sql
-- ❌ 对索引列使用函数
SELECT * FROM users WHERE UPPER(name) = 'ALICE';

-- ❌ 隐式类型转换
SELECT * FROM users WHERE id = '1';  -- id 是 INT 类型

-- ❌ LIKE 以通配符开头
SELECT * FROM users WHERE name LIKE '%Alice';

-- ✅ 这样可以使用索引
SELECT * FROM users WHERE name LIKE 'Alice%';
```

## 慢查询分析

使用 `EXPLAIN` 分析查询计划：

```sql
EXPLAIN SELECT * FROM users WHERE email = 'test@example.com'\G
```

关注 `type` 字段（性能从好到坏）：

| type | 说明 |
|------|------|
| `system` | 表只有一行 |
| `const` | 主键或唯一索引等值查询 |
| `ref` | 非唯一索引等值查询 |
| `range` | 索引范围查询 |
| `index` | 全索引扫描 |
| `ALL` | **全表扫描，需要优化！** |

## 索引设计原则

1. **选择性高**的列适合建索引（如 email，不适合 status）
2. **频繁查询**的列建索引，避免过多索引影响写入性能
3. 联合索引中，**选择性高**的列放在前面
4. 定期使用 `pt-query-digest` 分析慢查询日志
