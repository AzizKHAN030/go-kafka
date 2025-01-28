package utils

import (
	"time"

	"github.com/azizkhan030/go-kafka/registration/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int) (string, error) {
	jwtSecret := []byte(config.LoadConfig().JWTSecret)

	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	jwtSecret := []byte(config.LoadConfig().JWTSecret)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
