package security

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}
