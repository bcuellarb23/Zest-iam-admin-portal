package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Number used once this is a security check to avoid attack replication
func GenerateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("Couldn't generate nonce.")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

func GenerateState() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("Couldn't generate state.")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}
