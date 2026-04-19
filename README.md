# camelia's Blog

一个 MkDocs Material 风格的技术博客，Go 后端 + Vue 3 前端。

## 技术栈

**后端**
- Go + Gin（HTTP 框架）
- goldmark（Markdown 渲染）
- YAML frontmatter（文章元数据）

**前端**
- Vue 3 + Vite
- Tailwind CSS（Catppuccin Mocha 暗色主题）
- highlight.js（代码高亮）
- Vue Router（SPA 路由）

## 快速开始

### 开发模式

**终端 1 — 启动后端：**
```bash
cd backend
go run main.go
# 后端运行在 http://localhost:8080
```

**终端 2 — 启动前端：**
```bash
cd frontend
npm install  # 首次运行
npm run dev
# 前端运行在 http://localhost:5173
```

### 生产构建

```bash
# 构建前端并输出到 backend/dist
cd frontend && npm run build

# 启动后端（同时服务前端静态文件）
cd backend && go run main.go
# 访问 http://localhost:8080
```

## 添加文章

在 `backend/content/<分类>/` 目录下创建 `.md` 文件：

```markdown
---
title: "文章标题"
date: "2024-01-01"
tags: ["标签1", "标签2"]
excerpt: "文章摘要，显示在列表页。"
author: "Mayo"
---

# 正文标题

文章正文...
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/categories` | 获取所有分类 |
| GET | `/api/posts` | 获取所有文章 |
| GET | `/api/posts?category=xxx` | 按分类筛选文章 |
| GET | `/api/posts/:category/:slug` | 获取单篇文章（含 HTML） |

## 项目结构

```
blog/
├── backend/
│   ├── main.go              # 入口
│   ├── handlers/api.go      # API 处理器
│   ├── models/post.go       # 数据模型
│   └── content/             # Markdown 文章（新增目录即新增分类）
│       ├── algorithms/      # 算法
│       ├── database/        # 数据库
│       ├── network/         # 网络
│       ├── os/              # 操作系统
│       ├── go/              # Go 语言学习
│       └── mind/            # 生活的点点滴滴
└── frontend/
    ├── src/
    │   ├── components/
    │   │   ├── AppNavbar.vue    # 顶部导航（MkDocs 风格）
    │   │   └── PostCard.vue     # 文章卡片
    │   └── views/
    │       ├── HomeView.vue     # 首页
    │       ├── CategoryView.vue # 分类页
    │       └── PostView.vue     # 文章页（含 TOC）
    └── tailwind.config.js   # Catppuccin Mocha 主题
```
