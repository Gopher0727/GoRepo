package api

import (
	"github.com/gin-gonic/gin"

	"github.com/Gopher0727/GoRepo/backend/handlers"
	"github.com/Gopher0727/GoRepo/backend/middleware"
)

// RegisterRoutes 为 Gin Engine 注册所有业务路由
func RegisterRoutes(r *gin.Engine) {
	// 全局基础路由
	r.GET("/health", handlers.Health)

	api := r.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", handlers.AuthRegister)
		auth.POST("/login", handlers.AuthLogin)
		auth.GET("/rehash-check", handlers.AuthRehashCheck)
	}

	// 需要鉴权的受保护路由
	secure := v1.Group("/secure", middleware.Auth())
	{
		secure.POST("/encrypt", handlers.EncryptData)
		secure.POST("/decrypt", handlers.DecryptData)
	}

	// 示例：密码条目 CRUD（后续实现 handlers）
	entries := v1.Group("/entries", middleware.Auth())
	{
		entries.GET("", handlers.EntryList)          // 列表
		entries.POST("", handlers.EntryCreate)       // 创建
		entries.GET("/:id", handlers.EntryDetail)    // 详情
		entries.PUT("/:id", handlers.EntryUpdate)    // 全量更新
		entries.DELETE("/:id", handlers.EntryDelete) // 删除
	}
}
