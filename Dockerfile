# ============================================================
# Stage 1: 构建前端（Node.js）
# ============================================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
# 覆盖 outDir，在容器内输出到 ./dist
RUN npx vite build --outDir ./dist

# ============================================================
# Stage 2: 构建后端（Go）
# ============================================================
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 先复制 go.mod/go.sum，利用 Docker 层缓存加速依赖下载
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

# CGO_ENABLED=0 → 静态链接，可在 Alpine 运行
# -ldflags="-w -s" → 去除调试信息，缩小二进制体积
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o blog .

# ============================================================
# Stage 3: 最终镜像（极简 Alpine）
# ============================================================
FROM alpine:3.19

# 时区 + HTTPS 证书
RUN apk add --no-cache tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

WORKDIR /app

# 从各阶段复制产物
COPY --from=backend-builder /app/blog          ./blog
COPY --from=frontend-builder /app/dist         ./dist
COPY backend/content                           ./content

EXPOSE 8080

ENV GIN_MODE=release

CMD ["./blog"]
