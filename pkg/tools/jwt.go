package tools

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateToken(uid string, expiry int, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    uid,
		"exp":        time.Now().Add(time.Hour * time.Duration(24*expiry)).Unix(), // 过期时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}
