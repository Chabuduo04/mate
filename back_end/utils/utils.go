package utils

import (
	"fmt"
	"time"

	"github.com/Chabuduo04/mate/back_end/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"name": username,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := config.GetConfig().JWT.Secret
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT secret not configured")
	}
	return token.SignedString([]byte(jwtSecret))
}

func ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return config.GetConfig().JWT.Secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
