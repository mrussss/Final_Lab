// Package task 是任务领域模块，当前仅提供可编译的学习占位。
package task

import "time"

// Task 的模型设计留给 B1。
// TODO(B1): 根据学习手册定义 id、title、status、created_at、updated_at。Task、ID、Title

// TODO(B1): 自行设计字段约束、GORM tags 和请求/响应 DTO。
type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Status    string    `json:"status" gorm:"size:16;not null;default:pending"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskEvent 的关联事件模型留给 B2。
// TODO(B2): 定义 id、task_id、action、created_at 及关联/删除策略。
type TaskEvent struct{}
