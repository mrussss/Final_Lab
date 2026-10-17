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

func TestTaskHTTPIntegration(t *testing.T) {
	if os.Getenv("LAB_INTEGRATION_TEST") != "1" {
		t.Skip("integration test disabled")
	}

	// 1. 连接真实的实验数据库
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
		t.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	// 2. 使用真正的 Repository 和 Service
	repo := task.NewRepository(db)
	service := task.NewService(repo, nil)
	handler := task.NewHandler(service)

	router := NewRouter(
		filepath.Join("..", "..", "web", "index.html"),
		handler,
	)

	// 3. 定义测试中重复使用的 HTTP 请求函数
	send := func(method, path, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(
			method, path, strings.NewReader(body),
		)
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		return resp
	}

	// 4. POST：成功创建
	resp := send(
		http.MethodPost,
		"/tasks",
		`{"title":"B6 HTTP test"}`,
	)

	if resp.Code != http.StatusCreated {
		t.Fatalf("POST: status=%d, body=%s",
			resp.Code, resp.Body.String())
	}

	var created task.Task
	if err := json.Unmarshal(resp.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("POST did not return a valid task ID")
	}

	// 避免测试失败后留下自己的实验数据
	t.Cleanup(func() {
		result := db.Where("id = ?", created.ID).
			Delete(&task.Task{})
		if result.Error != nil {
			t.Errorf("cleanup task: %v", result.Error)
		}
	})

	// 5. GET：成功查询刚创建的任务
	path := fmt.Sprintf("/tasks/%d", created.ID)
	resp = send(http.MethodGet, path, "")

	if resp.Code != http.StatusOK {
		t.Fatalf("GET: status=%d", resp.Code)
	}

	var got task.Task
	if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}

	if got.ID != created.ID || got.Title != "B6 HTTP test" {
		t.Fatalf("unexpected task: %+v", got)
	}

	// 6. POST：非法标题
	resp = send(
		http.MethodPost,
		"/tasks",
		`{"title":"   "}`,
	)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	// 7. DELETE：删除刚创建的任务
	resp = send(http.MethodDelete, path, "")

	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.Code)
	}

	// 8. GET：已删除任务应当返回 404
	resp = send(http.MethodGet, path, "")

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}

	t.Log("PASS: create, get, validation, delete, not found")
}
