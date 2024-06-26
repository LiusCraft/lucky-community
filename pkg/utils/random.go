package utils

import (
	"math/rand"
	"time"
)

func GenerateCode(length int) string {

	rand.Seed(time.Now().UnixNano())

	charSet := "0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSet))
		code[i] = charSet[randomIndex]
	}
	return string(code)
}
