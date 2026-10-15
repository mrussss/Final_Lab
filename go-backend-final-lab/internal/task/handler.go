package task

import (
	"errors"
	"net/http"
	"strconv"

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
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"), 10, strconv.IntSize,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "invalid id"})
		return
	}

	task, err := h.service.GetTask(
		c.Request.Context(), uint(id),
	)

	switch {
	case errors.Is(err, ErrTaskNotFound):
		c.JSON(http.StatusNotFound,
			gin.H{"error": "task not found"})
	case err != nil:
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "internal_error"})
	default:
		c.JSON(http.StatusOK, task)
	}
}
func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	pageText := c.DefaultQuery("page", "1")
	pageSizeText := c.DefaultQuery("page_size", "10")
	page, err := strconv.Atoi(pageText)
	pageSize, err := strconv.Atoi(pageSizeText)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination"})
		return
	}
	items, err := h.service.ListTasks(
		c.Request.Context(),
		status,
		page,
		pageSize,
	)
	switch {
	case errors.Is(err, ErrInvalidStatus),
		errors.Is(err, ErrInvalidPagination):
		c.JSON(http.StatusNotFound,
			gin.H{"error": "task not found"})
	case err != nil:
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "internal_error"})
	default:

		c.JSON(http.StatusOK, gin.H{
			"items":     items,
			"page":      page,
			"page_size": pageSize,
		})
	}

}
func (*Handler) Update(c *gin.Context) { notImplemented(c) }
func (*Handler) Delete(c *gin.Context) { notImplemented(c) }

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}
