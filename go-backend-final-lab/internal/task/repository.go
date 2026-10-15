package task

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Repository 是持久化层占位，当前没有数据库连接与读写。
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

// TODO(B1): 根据实际调用需要设计方法，接入 GORM/MySQL，处理数据库错误。
// TODO(B1): 实现筛选、分页和删除策略，不在 Handler 中编写 SQL。
// TODO(B2): 配合 Service 实现事务关联写入与回滚验证。
type Repository struct {
	db *gorm.DB
}
