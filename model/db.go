package model

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dbConn, err := gorm.Open(sqlite.Open(os.Getenv("SQLITE_DB_PATH")), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db = dbConn
	db.AutoMigrate(&Voter{}, &Vote{})
}
