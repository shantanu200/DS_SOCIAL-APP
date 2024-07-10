package database

import (
	"fmt"
	"log"
	"notification/cmd/socket/models"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	dbURL := os.Getenv("DSN")

	start := time.Now()

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Notification Service Postgres  connected ", time.Since(start))

	db.AutoMigrate(&models.Notification{})
	db.AutoMigrate(&models.UserRelation{})

	DB = db
}
