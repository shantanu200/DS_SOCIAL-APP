package models

import (
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
	ProfileImage string    `json:"profileImage"`
}

type UserRelation struct {
	FolloweeId   int64     `json:"followeeid" gorm:"uniqueindex:user_relation_idx,constraint:onupdate:cascade,ondelete:set null;"`
	Followee     User      `json:"followee" gorm:"foreignKey:FolloweeId"`
	FollowerId   int64     `json:"followerId" gorm:"uniqueIndex:user_relation_idx,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Follower     User      `json:"follower" gorm:"foreignKey:FollowerId"`
	CreationTime time.Time `json:"createdAt"`
}
