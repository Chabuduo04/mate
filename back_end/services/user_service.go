package services

import (
	"errors"

	"github.com/Chabuduo04/mate/back_end/constants"
	"github.com/Chabuduo04/mate/back_end/dao"
	"github.com/Chabuduo04/mate/back_end/db"
	"github.com/Chabuduo04/mate/back_end/dto/respond"
	"github.com/Chabuduo04/mate/back_end/models"
	"github.com/Chabuduo04/mate/back_end/utils"
	"github.com/Chabuduo04/mate/back_end/zlog"
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

func (s *UserService) Register(username, password string) (string, *respond.RegisterRespond, int) {
	conn := db.GetConn()
	if _, err := dao.GetUserByUsername(conn, username); err == nil {
		return "用户已存在", nil, -2
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return constants.SYSTEM_ERROR, nil, -1
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return constants.SYSTEM_ERROR, nil, -1
	}
	u := &models.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashed),
	}
	if err := dao.CreateUser(conn, u); err != nil {
		zlog.Error("failed to create user: " + err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	return "注册成功", &respond.RegisterRespond{
		Id:       u.ID,
		Username: u.Username,
	}, 0
}

func (s *UserService) Authenticate(username, password string) (string, *respond.LoginRespond, int) {
	conn := db.GetConn()
	user, err := dao.GetUserByUsername(conn, username)
	if err != nil {
		return "用户不存在", nil, -2
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "密码错误", nil, -2
	}

	// create JWT
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return "生成Token失败", nil, -1
	}

	loginRsp := &respond.LoginRespond{
		Token:    token,
		Id:       user.ID,
		Username: user.Username,
	}
	return "登录成功", loginRsp, 0
}

func (s *UserService) GetByID(id string) (*models.User, bool) {
	conn := db.GetConn()
	u, err := dao.GetUserByID(conn, id)
	if err != nil {
		return nil, false
	}
	return u, true
}
