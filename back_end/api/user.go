package api

import (
	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/Chabuduo04/mate/back_end/utils"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	Service services.UserService
}

func NewUserApi(service services.UserService) UserApi {
	return UserApi{Service: service}
}

func (u *UserApi) UserRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := u.Service.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"id": user.ID, "username": user.Username})
}

func (u *UserApi) UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := u.Service.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// create JWT
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token, "id": user.ID, "username": user.Username})
}
