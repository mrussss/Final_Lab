package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"go-backend-final-lab/internal/task"
)

// 基础测试只验证 HTTP 脚手架；业务测试由学习者在 B1–B6 编写。
func TestHealthz(t *testing.T) {
	router := NewRouter(filepath.Join("..", "..", "web", "index.html"), task.NewHandler(nil))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("状态码 = %d，预期 200", response.Code)
	}
	if !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("预期 JSON，实际 Content-Type = %q", response.Header().Get("Content-Type"))
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("status = %q，预期 ok", body.Status)
	}
}

func TestStaticPage(t *testing.T) {
	router := NewRouter(filepath.Join("..", "..", "web", "index.html"), task.NewHandler(nil))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("状态码 = %d，预期 200", response.Code)
	}
	if !strings.HasPrefix(response.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("预期 HTML，实际 Content-Type = %q", response.Header().Get("Content-Type"))
	}
	if !strings.Contains(response.Body.String(), "Go Backend Final Lab") {
		t.Fatal("未返回实验静态页面")
	}
}
