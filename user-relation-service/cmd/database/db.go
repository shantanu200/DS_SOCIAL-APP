package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"user-relation/cmd/api/models"

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

	fmt.Println("User-Relation Service Postgres  connected ", time.Since(start))

    db.AutoMigrate(&models.User{},&models.UserRelation{})

	DB = db
}
