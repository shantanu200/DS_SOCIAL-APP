package models

import (
	pq "github.com/lib/pq"
	"time"
)

type User struct {
	Id           int       `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	CreationTime time.Time `json:"creationTime"`
	LastLogin    time.Time `json:"lastLogin"`
	IsHotUser    bool      `json:"isHotUser"`
	ProfileImage string `json:"profileImage"`
}

type Tweet struct {
	Id           int            `json:"id" gorm:"primaryKey"`
	UserId       int            `json:"userId" form:"userId"`
	User         User           `json:"user" gorm:"foreignKey:UserId"`
	CreationTime time.Time      `json:"createdAt"`
	Content      string         `json:"content" form:"content"`
	TotalLikes   int            `json:"totalLikes"`
	MediaFiles   pq.StringArray `json:"mediaFiles" gorm:"type:text[]" form:"mediaFiles"`
}