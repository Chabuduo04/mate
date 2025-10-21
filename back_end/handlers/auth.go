package handlers

import (
	"net/http"
	"time"

	"github.com/Chabuduo04/mate/back_end/config"
	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, svc *services.Services) {
	rg.POST("/register", makeRegisterHandler(svc))
	rg.POST("/login", makeLoginHandler(svc))
}

func makeRegisterHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Username string `json:"username" binding:"required,min=3"`
			Password string `json:"password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// type assertion to our UserService if present
		userSvc := svc.UserService
		if userSvc == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user service not available"})
			return
		}

		user, err := userSvc.Register(body.Username, body.Password)
		if err != nil {
			if err == services.ErrUserExists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// do not return password
		c.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username})
	}
}

func makeLoginHandler(svc *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userSvc := svc.UserService
		if userSvc == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user service not available"})
			return
		}

		user, err := userSvc.Authenticate(body.Username, body.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// create JWT
		jwtSecret := config.GetConfig().JWT.Secret
		if jwtSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt secret not configured"})
			return
		}

		claims := jwt.MapClaims{
			"sub":  user.ID,
			"name": user.Username,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
			return
		}

		// Optionally store session mapping if needed for revocation (skipped)

		c.JSON(http.StatusOK, gin.H{"token": signed, "user_id": user.ID, "username": user.Username})
	}
}
