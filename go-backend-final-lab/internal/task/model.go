package task

import "time"

type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Status    string    `json:"status" gorm:"size:16;not null;default:pending"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskEvent 的关联事件模型留给 B2。
// TODO(B2): 定义 id、task_id、action、created_at 及关联/删除策略。
type TaskEvent struct {
	ID        uint   `gorm:"primaryKey"`
	TaskID    uint   `gorm:"not null;index"`
	Action    string `gorm:"size:16;not null"`
	CreatedAt time.Time
}
