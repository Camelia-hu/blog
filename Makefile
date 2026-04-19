.PHONY: dev-backend dev-frontend install build run

# 安装前端依赖
install:
	cd frontend && npm install

# 开发模式：后端
dev-backend:
	cd backend && go run main.go

# 开发模式：前端（需先运行 dev-backend）
dev-frontend:
	cd frontend && npm run dev

# 构建生产版本
build:
	cd frontend && npm run build
	cd backend && go build -o ../bin/blog .

# 生产运行（需先 build）
run:
	./bin/blog

# 下载 Go 依赖
tidy:
	cd backend && go mod tidy
