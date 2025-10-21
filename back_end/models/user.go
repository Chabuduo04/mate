package models

type User struct {
    ID       string `json:"id" gorm:"primaryKey;type:varchar(36)"`
    Username string `json:"username" gorm:"uniqueIndex;size:128"`
    Password string `json:"password,omitempty"` // stored as bcrypt hash
}
