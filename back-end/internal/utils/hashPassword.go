package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

func GenerateSalt() string {
	n := rand.Intn(100)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := make([]byte, n)
	for i := range salt {
		salt[i] = charset[rand.Intn(len(charset))]
	}
	return string(salt)
}

func GeneratePassword(inputPassword, salt string) string {
	data := []byte(inputPassword + salt)
	hashedPassword := sha256.Sum256(data)
	return hex.EncodeToString(hashedPassword[:])
}

func VerifyPassword(inputPassword, salt, storedPassword string) bool {
	inputHash := GeneratePassword(inputPassword, salt)
	return inputHash == storedPassword
}
