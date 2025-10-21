package dao

import (
    "github.com/Chabuduo04/mate/back_end/models"
    "gorm.io/gorm"
)

// GetUserByUsername returns the user with the given username or gorm.ErrRecordNotFound
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
    var u models.User
    if err := db.Where("username = ?", username).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

// CreateUser inserts the given user into the database
func CreateUser(db *gorm.DB, u *models.User) error {
    return db.Create(u).Error
}

// GetUserByID returns the user with the given id or gorm.ErrRecordNotFound
func GetUserByID(db *gorm.DB, id string) (*models.User, error) {
    var u models.User
    if err := db.First(&u, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &u, nil
}
