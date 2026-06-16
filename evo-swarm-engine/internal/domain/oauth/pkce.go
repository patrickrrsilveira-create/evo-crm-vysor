package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// GeneratePKCE gera um code_verifier criptograficamente seguro e seu respectivo code_challenge
// baseado no RFC 7636 (S256).
func GeneratePKCE() (verifier string, challenge string, err error) {
	// 1. Gerar verifier (string pseudo-aleatória forte de alta entropia)
	// O RFC pede entre 43 e 128 caracteres. Vamos usar 32 bytes de entropia pura.
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return "", "", fmt.Errorf("falha ao gerar bytes aleatórios para PKCE: %w", err)
	}

	// Codifica para Base64-URL sem padding, o que resulta em exatamente 43 caracteres.
	verifier = base64.RawURLEncoding.EncodeToString(verifierBytes)

	// 2. Gerar o challenge via SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(verifier))
	hashBytes := hasher.Sum(nil)

	// Codifica o hash em Base64-URL sem padding
	challenge = base64.RawURLEncoding.EncodeToString(hashBytes)

	return verifier, challenge, nil
}
