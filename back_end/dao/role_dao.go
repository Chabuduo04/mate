package dao

import (
	"github.com/Chabuduo04/mate/back_end/models"
	"gorm.io/gorm"
)

func InsertRolesToDB(db *gorm.DB, roles []*models.Role) {
	for _, role := range roles {
		db.Create(&role)
	}
}
