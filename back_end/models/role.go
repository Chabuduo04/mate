package models

type Role struct {
	ID          string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string `json:"name" gorm:"type:varchar(128)"`
	Description string `json:"description" gorm:"type:varchar(256)"`
	Prompt      string `json:"prompt" gorm:"type:text"`
	VoiceType   string `json:"voice_type,omitempty" gorm:"type:varchar(64)"`
	CreatedAt   int64  `json:"created_at" gorm:"autoCreateTime;not null;comment:创建时间"`
	UpdatedAt   int64  `json:"updated_at" gorm:"autoUpdateTime;not null;comment:更新时间"`
}
