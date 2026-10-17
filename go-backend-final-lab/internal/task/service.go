package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidTitle      = errors.New("title cannot be empty or length greater than 100")
	ErrInvalidStatus     = errors.New("invalid task status")
	ErrInvalidPagination = errors.New("invalid pagination")
)

type TaskCache interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
}

type Service struct {
	repo  *Repository
	cache TaskCache
}

func NewService(repo *Repository, taskCache TaskCache) *Service {
	return &Service{
		repo:  repo,
		cache: taskCache,
	}
}

func taskCacheKey(id uint) string {
	return fmt.Sprintf("final_lab:task:%d", id)
}

func (s *Service) CreateTask(
	ctx context.Context,
	title string,
	status string,
) (*Task, error) {

	title = strings.TrimSpace(title)
	if title == "" || utf8.RuneCountInString(title) > 100 {
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

	key := taskCacheKey(id)

	// 1. 优先读取 Redis
	if s.cache != nil {
		raw, hit, err := s.cache.Get(ctx, key)

		switch {
		case err != nil:
			log.Printf("[cache] ERROR GET %s: %v", key, err)

		case hit:
			var cached Task

			if err := json.Unmarshal([]byte(raw), &cached); err == nil && cached.ID == id {
				log.Printf("[cache] HIT %s", key)
				return &cached, nil
			}

			log.Printf("[cache] invalid JSON or task ID: %s", key)

		default:
			log.Printf("[cache] MISS %s", key)
		}
	}

	// 2. 缓存未命中或故障时，回源 MySQL
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Printf("[db] GET task id=%d", id)

	// 3. 将 MySQL 结果写回 Redis
	if s.cache != nil {
		data, err := json.Marshal(task)

		if err != nil {
			log.Printf("[cache] JSON encode failed: %v", err)
		} else {
			err = s.cache.Set(ctx, key, string(data))
			if err != nil {
				log.Printf("[cache] SET failed %s: %v", key, err)
			} else {
				log.Printf("[cache] SET %s", key)
			}
		}
	}

	return task, nil
}

func (s *Service) invalidateTaskCache(
	ctx context.Context,
	id uint,
) {
	if s.cache == nil {
		return
	}

	key := taskCacheKey(id)

	if err := s.cache.Delete(ctx, key); err != nil {
		log.Printf("[cache] DEL failed %s: %v", key, err)
		return
	}

	log.Printf("[cache] DEL %s", key)
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
		if value == "" || utf8.RuneCountInString(value) > 100 {
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

	updated, err := s.repo.UpdateFields(ctx, id, changes)
	if err != nil {
		return nil, err
	}

	s.invalidateTaskCache(ctx, id)

	return updated, nil
}

func (s *Service) DeleteTask(
	ctx context.Context,
	id uint,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.invalidateTaskCache(ctx, id)
	return nil
}
