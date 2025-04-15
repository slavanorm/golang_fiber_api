package database

import (
	"log"
	"rentincome/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// global
var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	err = DB.AutoMigrate(&model.User{}, &model.Estate{}, &model.Transaction{})
	if err != nil {
		panic("Failed to migrate database")
	}
}

func CloseDB() {
	// could be used with non-sqlite db
	conn, err := DB.DB()
	if err != nil {
		panic("Failed to get database connection")
	}
	log.Print("closed conn")
	conn.Close()
}
