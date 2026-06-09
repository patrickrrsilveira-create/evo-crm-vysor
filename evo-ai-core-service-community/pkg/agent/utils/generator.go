package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

func GenerateAPIKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return uuid.New().String()
	}
	return hex.EncodeToString(bytes)
}
