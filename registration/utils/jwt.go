package utils

import (
	"fmt"
	"time"

	"github.com/azizkhan030/go-kafka/registration/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int) (string, error) {
	jwtSecret := config.LoadConfig().JWTSecret

	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token:%v", jwtSecret)

	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return config.LoadConfig().JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
