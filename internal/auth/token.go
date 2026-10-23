package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// This structure can be modified to configure parameters
type TokenGenerator struct{}

func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

func (g *TokenGenerator) Generate() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate securre token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
