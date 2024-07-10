package models

import "time"

type User struct {
	Id           int       `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	Email        string    `json:"email" gorm:"uniqueIndex:user_email_unique_idx,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Password     string    `json:"password"`
	CreationTime time.Time `json:"creationTime"`
	LastLogin    time.Time `json:"lastLogin"`
	IsHotUser    bool      `json:"isHotUser"`
	ProfileImage string    `json:"profileImage"`
}
