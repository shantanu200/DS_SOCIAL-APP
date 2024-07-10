package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"tweet/cmd/api/models"

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

	fmt.Println("Tweet Service Postgres connected ", time.Since(start))

	db.AutoMigrate(&models.Tweet{})
	db.AutoMigrate(&models.FavouriteTweet{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Reply{})
	db.AutoMigrate(&models.Thread{})
	db.AutoMigrate(&models.FavouriteReply{})
	DB = db
}
