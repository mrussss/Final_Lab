package main

import (
	"fmt"
	"log"
	"os"

	"go-backend-final-lab/internal/app"
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

	// TODO：调用 sqlDB.Ping()
	err = sqlDB.Ping()
	// TODO：如果 Ping 返回错误，使用 log.Fatalf 处理
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	// TODO：如果成功，使用 log.Println 输出成功信息
	log.Println("数据库连接成功")

	if err := db.AutoMigrate(&task.Task{}); err != nil {
		log.Fatalf("数据库表迁移失败: %v", err)
	}
	log.Println("Task 数据表迁移完成")

	// 从项目根目录启动，静态文件路径相对于当前工作目录。
	router := app.NewRouter("web/index.html")
	log.Printf("HTTP 服务启动，端口 %s；/healthz 仅检查 HTTP 进程", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}
}
