package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabase2() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/db_store_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	return db
}
