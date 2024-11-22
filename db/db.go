package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB là biến toàn cục để lưu kết nối database
var DB *gorm.DB

// Connect thiết lập kết nối với cơ sở dữ liệu
func Connect() error {
	// Lấy thông tin kết nối từ biến môi trường
	dsn := getDSN()
	if dsn == "" {
		return fmt.Errorf("missing database configuration")
	}

	// Kết nối với PostgreSQL
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Hiển thị log query (tuỳ chỉnh mức độ log nếu cần)
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// getDSN lấy DSN từ biến môi trường
func getDSN() string {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	sslmode := os.Getenv("POSTGRES_SSL_MODE")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}
