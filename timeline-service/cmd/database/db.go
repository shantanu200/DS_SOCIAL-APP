package database

import (
	"log"
	"os"
	"time"
	"timeline/cmd/api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	dbURL := os.Getenv("DSN")

	start := time.Now()

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("TimeLine Service Postgres  connected ", time.Since(start))

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Tweet{})

	DB = db
}
