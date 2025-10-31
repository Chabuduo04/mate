package dao

import (
	"github.com/Chabuduo04/mate/back_end/models"
	"gorm.io/gorm"
)

func GetChatRecordsByUserAndRole(db *gorm.DB, userID string, roleID string) ([]models.ChatRecord, error) {
	var records []models.ChatRecord
	if err := db.Where("user_id = ? AND role_id = ?", userID, roleID).Order("created_at asc").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func CreateChatRecord(db *gorm.DB, record *models.ChatRecord) error {
	return db.Create(record).Error
}
