package task

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 是 HTTP 边界占位；尚未绑定 Service。
// TODO(B1): 解析和校验请求，将业务交给 Service，把结果映射为 HTTP 响应。
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) Create(c *gin.Context) {
	type CreateTaskRequest struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	created, err := h.service.CreateTask(
		c.Request.Context(), req.Title, req.Status,
	)
	switch {
	case errors.Is(err, ErrInvalidTitle),
		errors.Is(err, ErrInvalidStatus):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
	default:
		c.JSON(http.StatusCreated, created)
	}

}
func (*Handler) Get(c *gin.Context)    { notImplemented(c) }
func (*Handler) List(c *gin.Context)   { notImplemented(c) }
func (*Handler) Update(c *gin.Context) { notImplemented(c) }
func (*Handler) Delete(c *gin.Context) { notImplemented(c) }

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}
