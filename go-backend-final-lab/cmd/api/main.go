package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go-backend-final-lab/internal/app"
	"go-backend-final-lab/internal/cache"
	"go-backend-final-lab/internal/task"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "18080"
	}

	host := os.Getenv("LAB_MYSQL_HOST")
	user := os.Getenv("LAB_MYSQL_USER")
	password := os.Getenv("LAB_MYSQL_PASSWORD")
	dbName := os.Getenv("LAB_MYSQL_DATABASE")
	mysqlPort := os.Getenv("LAB_MYSQL_PORT")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, mysqlPort, dbName,
	)
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalf("GORM 初始化失败: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}
	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	log.Println("数据库连接成功")

	if err := db.AutoMigrate(
		&task.Task{},
		&task.TaskEvent{},
	); err != nil {
		log.Fatalf("数据库表迁移失败: %v", err)
	}
	log.Println("Task 数据表迁移完成")

	// 从项目根目录启动，静态文件路径相对于当前工作目录。

	redisAddr := os.Getenv("LAB_REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:16379"
	}

	// 本地实验使用 30 秒 TTL
	redisStore := cache.NewRedisStore(
		redisAddr,
		2*time.Minute,
	)
	defer redisStore.Close()

	// Redis 不可用时不终止服务
	if err := redisStore.Ping(context.Background()); err != nil {
		log.Printf("[cache] Redis unavailable at startup: %v", err)
	} else {
		log.Printf("[cache] Redis connected: %s", redisAddr)
	}

	repo := task.NewRepository(db)
	service := task.NewService(repo, redisStore)
	handler := task.NewHandler(service)

	router := app.NewRouter("web/index.html", handler)
	log.Printf("HTTP 服务启动，端口 %s；/healthz 仅检查 HTTP 进程", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}
}
