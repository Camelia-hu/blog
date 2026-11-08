---
title: "HTTP 协议深度解析"
date: "2024-02-01"
tags: ["网络", "HTTP", "HTTPS", "协议"]
excerpt: "从 HTTP/1.0 到 HTTP/3，深入理解超文本传输协议的演进历程、核心机制与 HTTPS 握手过程。"
author: "Mayo"
---

# HTTP 协议深度解析

HTTP（HyperText Transfer Protocol）是 Web 的基础协议，理解它对于构建高性能网络应用至关重要。

## HTTP 版本演进

### HTTP/1.0

- 每次请求建立新的 TCP 连接，请求完成后关闭
- 不支持持久连接，开销巨大
- 不支持虚拟主机（无 Host 头部）

### HTTP/1.1

相较于 1.0 的重要改进：

- **持久连接**：默认 `Connection: keep-alive`，复用 TCP 连接
- **管道化**（Pipelining）：可以连续发送多个请求，但存在队头阻塞问题
- **Host 头部**：支持虚拟主机
- **分块传输**：`Transfer-Encoding: chunked`

### HTTP/2

HTTP/2 是一次重大革新：

- **二进制分帧**：将消息分割为帧，更高效
- **多路复用**：单个 TCP 连接并发多个请求，解决 HTTP 层队头阻塞
- **头部压缩**：HPACK 算法，大幅减少头部体积
- **服务器推送**：主动向客户端推送资源

### HTTP/3

基于 QUIC 协议（UDP），从根本上解决了问题：

- **解决 TCP 队头阻塞**：每个流独立，丢包只影响单个流
- **0-RTT 连接建立**：缓存连接信息，大幅降低延迟
- **内置 TLS 1.3**：加密是 QUIC 的核心组成部分
- **连接迁移**：IP 变化后连接不中断（移动网络友好）

## HTTP 请求结构

```
GET /api/posts?page=1 HTTP/1.1
Host: example.com
User-Agent: Mozilla/5.0
Accept: application/json
Authorization: Bearer eyJhbGci...
```

## HTTP 响应结构

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Cache-Control: max-age=3600
X-Request-ID: abc123

{"posts": [...]}
```

## 常用状态码

| 状态码 | 含义 | 说明 |
|--------|------|------|
| 200 | OK | 请求成功 |
| 201 | Created | 资源创建成功 |
| 204 | No Content | 请求成功，无响应体 |
| 301 | Moved Permanently | 永久重定向 |
| 302 | Found | 临时重定向 |
| 304 | Not Modified | 资源未修改，使用缓存 |
| 400 | Bad Request | 请求格式错误 |
| 401 | Unauthorized | 未认证 |
| 403 | Forbidden | 无权限 |
| 404 | Not Found | 资源不存在 |
| 429 | Too Many Requests | 请求频率过高 |
| 500 | Internal Server Error | 服务器内部错误 |
| 502 | Bad Gateway | 网关错误 |
| 503 | Service Unavailable | 服务不可用 |

## HTTPS 握手过程（TLS 1.3）

TLS 1.3 大幅简化了握手流程：

```
Client                                Server
  |                                     |
  |--- ClientHello (支持的算法套件) ---->|
  |                                     |
  |<-- ServerHello + 证书 + Finished ---|
  |                                     |
  |--- Finished ----------------------->|
  |                                     |
  |<====== 加密通信开始 ===============>|
```

相比 TLS 1.2 减少了 1 个 RTT，0-RTT 模式下可以零延迟恢复连接。

## 缓存控制

```
# 强缓存（不请求服务器）
Cache-Control: max-age=31536000
Expires: Wed, 21 Oct 2025 07:28:00 GMT

# 协商缓存（请求服务器验证）
Cache-Control: no-cache
ETag: "abc123"
Last-Modified: Wed, 21 Oct 2023 07:28:00 GMT
```
