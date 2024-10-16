package tools

import "math/rand"

func RandomNumber(length int) string {
	letters := "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
