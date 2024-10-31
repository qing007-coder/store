package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt 加密
func Encrypt(plainText []byte, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 填充
	plainText = pad(plainText, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密
func Decrypt(cipherText, key string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去填充
	ciphertext = unpad(ciphertext)

	return ciphertext, nil
}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func unpad(src []byte) []byte {
	length := len(src)
	unpad := int(src[length-1])
	return src[:(length - unpad)]
}
