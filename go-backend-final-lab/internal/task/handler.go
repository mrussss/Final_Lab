package task

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination"})
		return
	}
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
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
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
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Title  *string `json:"title"`
		Status *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	updated, err := h.service.UpdateTask(
		c.Request.Context(), uint(id), req.Title, req.Status,
	)

	switch {
	case errors.Is(err, ErrInvalidTitle),
		errors.Is(err, ErrInvalidStatus),
		errors.Is(err, ErrEmptyUpdate):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrTaskNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
	default:
		c.JSON(http.StatusOK, updated)
	}
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.DeleteTask(c.Request.Context(), uint(id))

	switch {
	case errors.Is(err, ErrTaskNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
	default:
		c.Status(http.StatusNoContent)
	}
}
