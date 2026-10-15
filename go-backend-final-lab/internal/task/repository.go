package task

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrTaskNotFound = errors.New("task not found")

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, task *Task) error {
	result := r.db.WithContext(ctx).Create(task)
	return result.Error

}

func (r *Repository) GetByID(
	ctx context.Context,
	id uint,
) (*Task, error) {

	var task Task

	err := r.db.WithContext(ctx).
		First(&task, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("query task: %w", err)
	}

	return &task, nil
}

func (r *Repository) List(
	ctx context.Context,
	status string,
	page int,
	pageSize int,
) ([]Task, error) {

	tasks := make([]Task, 0)

	query := r.db.WithContext(ctx).Model(&Task{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	offset := (page - 1) * pageSize

	err := query.
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, nil
}

func (r *Repository) UpdateFields(
	ctx context.Context,
	id uint,
	changes map[string]any,
) (*Task, error) {
	result := r.db.WithContext(ctx).
		Model(&Task{}).
		Where("id = ?", id).
		Updates(changes)

	if result.Error != nil {
		return nil, fmt.Errorf("update task: %w", result.Error)
	}

	// 再次查询：既获得最新数据，也能识别不存在的 ID。
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(
	ctx context.Context,
	id uint,
) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&Task{})

	if result.Error != nil {
		return fmt.Errorf("delete task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// TODO(B2): 配合 Service 实现事务关联写入与回滚验证。
type Repository struct {
	db *gorm.DB
}
