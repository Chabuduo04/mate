package services

import (
	"errors"

	"github.com/Chabuduo04/mate/back_end/dao"
	"github.com/Chabuduo04/mate/back_end/db"
	"github.com/Chabuduo04/mate/back_end/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(username, password string) (*models.User, error) {
	conn := db.GetConn()
	if _, err := dao.GetUserByUsername(conn, username); err == nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashed),
	}
	if err := dao.CreateUser(conn, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Authenticate(username, password string) (*models.User, error) {
	conn := db.GetConn()
	u, err := dao.GetUserByUsername(conn, username)
	if err != nil {
		return nil, ErrInvalidCreds
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCreds
	}
	return u, nil
}

func (s *UserService) GetByID(id string) (*models.User, bool) {
	conn := db.GetConn()
	u, err := dao.GetUserByID(conn, id)
	if err != nil {
		return nil, false
	}
	return u, true
}
