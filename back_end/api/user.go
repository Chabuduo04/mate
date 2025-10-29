package api

import (
	"net/http"

	"github.com/Chabuduo04/mate/back_end/constants"
	"github.com/Chabuduo04/mate/back_end/dto/request"
	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/Chabuduo04/mate/back_end/zlog"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	Service services.UserService
}

func NewUserApi(service services.UserService) UserApi {
	return UserApi{Service: service}
}

func (u *UserApi) UserRegister(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}

	message, userInfo, ret := u.Service.Register(req.Username, req.Password)

	JsonBack(c, message, ret, userInfo)
}

func (u *UserApi) UserLogin(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}

	message, userInfo, ret := u.Service.Authenticate(req.Username, req.Password)

	JsonBack(c, message, ret, userInfo)
}
