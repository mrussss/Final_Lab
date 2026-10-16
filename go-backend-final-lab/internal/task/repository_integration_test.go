package task

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUpdateFieldsRollback(t *testing.T) {
	// 只有显式开启时才连接真实 MySQL。
	if os.Getenv("LAB_INTEGRATION_TEST") != "1" {
		t.Skip("MySQL integration test is disabled")
	}

	// 1. 连接 Final Lab 的 MySQL。
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("LAB_MYSQL_USER"),
		os.Getenv("LAB_MYSQL_PASSWORD"),
		os.Getenv("LAB_MYSQL_HOST"),
		os.Getenv("LAB_MYSQL_PORT"),
		os.Getenv("LAB_MYSQL_DATABASE"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("connect MySQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	if !db.Migrator().HasTable(&TaskEvent{}) {
		t.Fatal("task_events table does not exist")
	}

	repo := NewRepository(db)
	ctx := context.Background()

	// 2. 创建一条独立的测试任务。
	task := &Task{
		Title:  "B2 rollback integration test",
		Status: "pending",
	}

	if err := repo.Create(ctx, task); err != nil {
		t.Fatalf("create test task: %v", err)
	}

	// 测试结束后清理自己创建的数据。
	t.Cleanup(func() {
		cleanupCtx := context.Background()

		if err := db.WithContext(cleanupCtx).
			Where("task_id = ?", task.ID).
			Delete(&TaskEvent{}).Error; err != nil {
			t.Errorf("cleanup events: %v", err)
		}

		if err := db.WithContext(cleanupCtx).
			Where("id = ?", task.ID).
			Delete(&Task{}).Error; err != nil {
			t.Errorf("cleanup task: %v", err)
		}
	})

	// 3. 注册仅针对 TaskEvent 创建操作的失败注入。
	injectedErr := errors.New("injected event insert failure")
	const callbackName = "b2:test:reject_event_insert"

	err = db.Callback().Create().
		Before("gorm:create").
		Register(callbackName, func(tx *gorm.DB) {
			if tx.Statement != nil &&
				tx.Statement.Schema != nil &&
				tx.Statement.Schema.Table == "task_events" {
				tx.AddError(injectedErr)
			}
		})

	if err != nil {
		t.Fatalf("register callback: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Callback().Create().
			Remove(callbackName); err != nil {
			t.Errorf("remove callback: %v", err)
		}
	})

	// 4. 调用真正的 Repository 更新方法。
	updated, err := repo.UpdateFields(
		ctx,
		task.ID,
		map[string]any{"status": "done"},
	)

	// 必须收到我们故意注入的错误。
	if !errors.Is(err, injectedErr) {
		t.Fatalf("expected injected error, got: %v", err)
	}

	if updated != nil {
		t.Fatalf("expected nil result, got: %+v", updated)
	}

	// 5. 查询 MySQL，确认 UPDATE 已被回滚。
	got, err := repo.GetByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("read task after rollback: %v", err)
	}

	if got.Status != "pending" {
		t.Fatalf(
			"rollback failed: status = %s, want pending",
			got.Status,
		)
	}

	// 6. 确认没有产生对应的事件。
	var eventCount int64

	err = db.Model(&TaskEvent{}).
		Where("task_id = ?", task.ID).
		Count(&eventCount).Error

	if err != nil {
		t.Fatalf("count events: %v", err)
	}

	if eventCount != 0 {
		t.Fatalf(
			"rollback failed: %d events remain",
			eventCount,
		)
	}

	t.Log("PASS: Task update rolled back, no event inserted")
}
