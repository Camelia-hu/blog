package main

import (
	"blog/handlers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS 配置，开发时允许前端 dev server 访问
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}))

	// API 路由组
	api := r.Group("/api")
	{
		api.GET("/categories", handlers.GetCategories)
		api.GET("/posts", handlers.GetPosts)
		api.GET("/posts/:category/:slug", handlers.GetPost)
	}

	// 生产环境：服务前端静态文件
	r.Static("/assets", "./dist/assets")
	r.StaticFile("/favicon.svg", "./dist/favicon.svg")
	r.NoRoute(func(c *gin.Context) {
		// API 路由未匹配时返回 404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		// SPA 路由，返回 index.html
		c.File("./dist/index.html")
	})

	r.Run(":8080")
}
