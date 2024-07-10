package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"user/cmd/api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	dbURL := os.Getenv("DSN")

	start := time.Now()

	fmt.Printf("Trying to connect %s ",dbURL);
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("User Service Postgres  connected ", time.Since(start))

	db.AutoMigrate(&models.User{})

	DB = db
}
