// Package app 负责组装 HTTP 路由，不承担 Task 业务或基础设施连接。
package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-backend-final-lab/internal/task"
)

// NewRouter 注册基础服务与占位路由。indexPath 由启动入口指定。
// /healthz 不连接 MySQL 或 Redis，也不代表这些依赖已经就绪。
func NewRouter(indexPath string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/", func(c *gin.Context) {
		c.File(indexPath)
	})

	handler := &task.Handler{}
	router.POST("/tasks", handler.Create)
	router.GET("/tasks/:id", handler.Get)
	router.GET("/tasks", handler.List)
	router.PATCH("/tasks/:id", handler.Update)
	router.DELETE("/tasks/:id", handler.Delete)

	return router
}
