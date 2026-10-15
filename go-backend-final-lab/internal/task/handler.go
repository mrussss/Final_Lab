package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 是 HTTP 边界占位；尚未绑定 Service。
// TODO(B1): 解析和校验请求，将业务交给 Service，把结果映射为 HTTP 响应。
type Handler struct{}

func (*Handler) Create(c *gin.Context) {
	type CreateTaskRequest struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	var s Service
	Task, err := s.CreateTask(c.Request.Context(), req.Title, req.Status)
	if 


}
func (*Handler) Get(c *gin.Context)    { notImplemented(c) }
func (*Handler) List(c *gin.Context)   { notImplemented(c) }
func (*Handler) Update(c *gin.Context) { notImplemented(c) }
func (*Handler) Delete(c *gin.Context) { notImplemented(c) }

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}
