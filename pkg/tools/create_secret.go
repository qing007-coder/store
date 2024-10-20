package tools

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func GenerateSecret(length int) (string, error) {
	secret := make([]byte, length)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	base64Secret := base64.URLEncoding.EncodeToString(secret)
	base64Secret = strings.TrimRight(base64Secret, "=")
	return base64Secret, nil
}
