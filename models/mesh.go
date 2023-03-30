package models

import (
	"encoding/base64"
	"fmt"
	"math/rand"
)

func GenerateNonce(length int) (string, error) {
	if length <= 0 || length > 2048 {
		return "", fmt.Errorf("Length must be bewteen 1 and 2048")
	}

	nonceBytes := make([]byte, length)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}
