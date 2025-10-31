package models

import "time"

type ChatRecord struct {
	ID        int64     `json:"id" gorm:"primaryKey;comment:自增ID"`
	UserID    string    `json:"user_id" gorm:"index;not null;type:varchar(36)"`
	RoleID    string    `json:"role_id" gorm:"index;not null;type:varchar(36)"`
	Message   string    `json:"message" gorm:"type:text"`
	Response  string    `json:"response" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;not null;comment:创建时间"`
}
