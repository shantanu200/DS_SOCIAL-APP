package models

import "time"

type Notification struct {
	Id           int       `json:"id" gorm:"primaryKey"`
	Message      string    `json:"message"`
	UserId       int       `json:"userid" gorm:"uniqueindex:notification_u_idx,constraint:onupdate:cascade,ondelete:set null;"`
	User         User      `json:"user"`
	Type         string    `json:"type" gorm:"uniqueindex:notification_u_idx,constraint:onupdate:cascade,ondelete:set null;"`
	Read         bool      `json:"read"`
	CreationTime time.Time `json:"createdAt"`
	AuthorId     int       `json:"authorid" gorm:"uniqueindex:notification_u_idx,constraint:onupdate:cascade,ondelete:set null;"`
	Author       User      `json:"author"`
	PostId       int      `json:"postid" gorm:"uniqueindex:notification_u_idx,constraint:onupdate:cascade,ondelete:set null;"`
}
