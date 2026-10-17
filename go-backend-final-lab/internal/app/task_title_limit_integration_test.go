package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-backend-final-lab/internal/task"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestTaskTitleLimitIntegration(t *testing.T) {
	if os.Getenv("LAB_INTEGRATION_TEST") != "1" {
		t.Skip("integration test disabled")
	}

	// 1. 连接真实 MySQL
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

	// 2. 复用项目的真实业务层
	repo := task.NewRepository(db)
	service := task.NewService(repo, nil)
	handler := task.NewHandler(service)

	router := NewRouter(
		filepath.Join("..", "..", "web", "index.html"),
		handler,
	)

	// 构造 JSON 请求并直接交给 Gin Router
	send := func(method, path, title string) *httptest.ResponseRecorder {
		body, err := json.Marshal(map[string]string{
			"title": title,
		})
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(
			method,
			path,
			strings.NewReader(string(body)),
		)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		return resp
	}

	// 3. POST：四种边界情况
	cases := []struct {
		name       string
		title      string
		wantStatus int
	}{
		{
			name:       "English 100",
			title:      strings.Repeat("a", 100),
			wantStatus: http.StatusCreated,
		},
		{
			name:       "English 101",
			title:      strings.Repeat("a", 101),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Chinese 100",
			title:      strings.Repeat("中", 100),
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Chinese 101",
			title:      strings.Repeat("中", 101),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run("POST "+tc.name, func(t *testing.T) {
			resp := send(http.MethodPost, "/tasks", tc.title)

			if resp.Code != tc.wantStatus {
				t.Fatalf(
					"status=%d, want=%d, body=%s",
					resp.Code,
					tc.wantStatus,
					resp.Body.String(),
				)
			}

			// 创建成功时检查响应，并清理测试数据
			if resp.Code == http.StatusCreated {
				var created task.Task

				if err := json.Unmarshal(
					resp.Body.Bytes(), &created,
				); err != nil {
					t.Fatal(err)
				}

				if created.ID == 0 || created.Title != tc.title {
					t.Fatalf("unexpected task: %+v", created)
				}

				id := created.ID

				t.Cleanup(func() {
					if err := db.Where("id = ?", id).
						Delete(&task.Task{}).Error; err != nil {
						t.Errorf("cleanup task: %v", err)
					}
				})
			}
		})
	}

	// 4. PATCH：检查边界与数据库原值
	t.Run("PATCH boundary", func(t *testing.T) {

		// 创建专门的测试任务
		seed := &task.Task{
			Title:  "Original title",
			Status: "pending",
		}

		if err := db.Create(seed).Error; err != nil {
			t.Fatal(err)
		}

		id := seed.ID
		path := fmt.Sprintf("/tasks/%d", id)

		// 测试结束后清理事件与任务
		t.Cleanup(func() {
			if err := db.Where("task_id = ?", id).
				Delete(&task.TaskEvent{}).Error; err != nil {
				t.Errorf("cleanup events: %v", err)
			}

			if err := db.Where("id = ?", id).
				Delete(&task.Task{}).Error; err != nil {
				t.Errorf("cleanup task: %v", err)
			}
		})

		// PATCH：100 个中文字符，应成功
		validTitle := strings.Repeat("中", 100)

		resp := send(http.MethodPatch, path, validTitle)

		if resp.Code != http.StatusOK {
			t.Fatalf(
				"PATCH 100: status=%d, body=%s",
				resp.Code, resp.Body.String(),
			)
		}

		var afterValid task.Task
		if err := db.First(&afterValid, id).Error; err != nil {
			t.Fatal(err)
		}

		if afterValid.Title != validTitle {
			t.Fatal("PATCH 100 did not persist title")
		}

		// 记录当前 TaskEvent 数量
		var eventsBefore int64

		if err := db.Model(&task.TaskEvent{}).
			Where("task_id = ?", id).
			Count(&eventsBefore).Error; err != nil {
			t.Fatal(err)
		}

		// PATCH：101 个中文字符，应被拒绝
		invalidTitle := strings.Repeat("中", 101)

		resp = send(http.MethodPatch, path, invalidTitle)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf(
				"PATCH 101: status=%d, body=%s",
				resp.Code, resp.Body.String(),
			)
		}

		// 再次查询真实数据库，确认标题没变
		var afterInvalid task.Task

		if err := db.First(&afterInvalid, id).Error; err != nil {
			t.Fatal(err)
		}

		if afterInvalid.Title != validTitle {
			t.Fatal("invalid PATCH modified database title")
		}

		// 校验失败不应该产生新的 TaskEvent
		var eventsAfter int64

		if err := db.Model(&task.TaskEvent{}).
			Where("task_id = ?", id).
			Count(&eventsAfter).Error; err != nil {
			t.Fatal(err)
		}

		if eventsAfter != eventsBefore {
			t.Fatalf(
				"invalid PATCH created event: before=%d after=%d",
				eventsBefore,
				eventsAfter,
			)
		}

		t.Log("PATCH 100/101 boundary verified")
	})
}
