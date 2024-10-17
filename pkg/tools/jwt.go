package tools

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"store/pkg/config"
	"time"
)

func CreateToken(uid string, conf *config.GlobalConfig) (string, error) {
	expiry := conf.JWT.Expiry
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    uid,
		"exp":        time.Now().Add(time.Hour * time.Duration(24*expiry)).Unix(), // 过期时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(conf.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, conf *config.GlobalConfig) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.JWT.SecretKey), nil
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
