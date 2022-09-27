package urlshortener

import (
	"errors"
	"math/rand"
	"strings"
)

var ErrKeyGenerationFailed = errors.New("key_generation_failed")

func GenerateKey() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	const keyLength = 5
	var keyBuilder strings.Builder
	for i := 0; i < keyLength; i++ {
		char := alphabet[rand.Intn(len(alphabet))]
		_ = keyBuilder.WriteByte(char) // string builder never fails writes
	}
	return keyBuilder.String()
}
