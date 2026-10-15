package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidTitle      = errors.New("title cannot be empty")
	ErrInvalidStatus     = errors.New("invalid task status")
	ErrInvalidPagination = errors.New("invalid pagination")
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

func (s *Service) GetTask(
	ctx context.Context,
	id uint,
) (*Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListTasks(
	ctx context.Context,
	status string,
	page int,
	pageSize int,
) ([]Task, error) {

	if status != "" &&
		status != "pending" &&
		status != "done" {
		return nil, ErrInvalidStatus
	}

	if page < 1 || page > 10000 ||
		pageSize < 1 || pageSize > 100 {
		return nil, ErrInvalidPagination
	}

	return s.repo.List(ctx, status, page, pageSize)
}

// TODO(B2): 实现 PATCH 更新和事件插入的事务边界。
// TODO(B3): 实现缓存读取、回源以及提交后的缓存失效。
var ErrEmptyUpdate = errors.New("no fields to update")

func (s *Service) UpdateTask(
	ctx context.Context,
	id uint,
	title *string,
	status *string,
) (*Task, error) {
	if title == nil && status == nil {
		return nil, ErrEmptyUpdate
	}

	changes := make(map[string]any)

	if title != nil {
		value := strings.TrimSpace(*title)
		if value == "" {
			return nil, ErrInvalidTitle
		}
		changes["title"] = value
	}

	if status != nil {
		if *status != "pending" && *status != "done" {
			return nil, ErrInvalidStatus
		}
		changes["status"] = *status
	}

	return s.repo.UpdateFields(ctx, id, changes)
}

func (s *Service) DeleteTask(
	ctx context.Context,
	id uint,
) error {
	return s.repo.Delete(ctx, id)
}
