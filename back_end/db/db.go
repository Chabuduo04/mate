package db

import (
	"fmt"
	"log"

	"github.com/Chabuduo04/mate/back_end/config"
	"github.com/Chabuduo04/mate/back_end/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

func InitDB() {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	gormDB = db
}

func GetConn() *gorm.DB {
	return gormDB
}
