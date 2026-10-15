package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Service 预留业务职责，不依赖 *gin.Context。
var (
	ErrInvalidTitle  = errors.New("title cannot be empty")
	ErrInvalidStatus = errors.New("invalid task status")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(
	ctx context.Context,
	title string,
	status string,
) (*Task, error) {

	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrInvalidTitle
	}

	if status == "" {
		status = "pending"
	}

	if status != "pending" && status != "done" {
		return nil, ErrInvalidStatus
	}

	task := &Task{
		Title:  title,
		Status: status,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

// TODO(B1): 设计业务方法、context.Context 参数与 Repository 依赖。
// TODO(B2): 实现 PATCH 更新和事件插入的事务边界。
// TODO(B3): 实现缓存读取、回源以及提交后的缓存失效。
// 当前没有业务方法、数据库操作或缓存操作。
