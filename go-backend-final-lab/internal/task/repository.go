package task

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	var updated Task

	err := r.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {

			// 1. 查询 Task 并锁定目标记录
			err := tx.Clauses(
				clause.Locking{Strength: "UPDATE"},
			).First(&updated, id).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTaskNotFound
			}
			if err != nil {
				return fmt.Errorf("find task: %w", err)
			}

			// 2. 更新 Task
			if err := tx.Model(&updated).
				Updates(changes).Error; err != nil {
				return fmt.Errorf("update task: %w", err)
			}

			// 3. 插入操作事件
			event := TaskEvent{
				TaskID: id,
				Action: "update",
			}

			if err := tx.Create(&event).Error; err != nil {
				return fmt.Errorf("insert event: %w", err)
			}

			// 4. 在同一个事务内读取更新结果
			if err := tx.First(&updated, id).Error; err != nil {
				return fmt.Errorf("reload task: %w", err)
			}

			return nil
		},
	)

	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf(
			"update task transaction: %w", err,
		)
	}

	return &updated, nil
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
