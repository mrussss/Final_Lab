package task

import (
	"context"

	"gorm.io/gorm"
)

// Repository 是持久化层占位，当前没有数据库连接与读写。

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, task *Task) error {
	result := r.db.WithContext(ctx).Create(task)
	return result.Error

}

// TODO(B1): 根据实际调用需要设计方法，接入 GORM/MySQL，处理数据库错误。
// TODO(B1): 实现筛选、分页和删除策略，不在 Handler 中编写 SQL。
// TODO(B2): 配合 Service 实现事务关联写入与回滚验证。
type Repository struct {
	db *gorm.DB
}
