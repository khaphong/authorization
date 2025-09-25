package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime     = 1
	argonMemory   = 64 * 1024
	argonThreads  = 4
	argonKeyLen   = 32
	saltLen       = 16
)

// HashPassword hashes a password using argon2id
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Combine salt and hash
	combined := append(salt, hash...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	combined, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}

	if len(combined) != saltLen+argonKeyLen {
		return false, fmt.Errorf("invalid hash format")
	}

	salt := combined[:saltLen]
	hash := combined[saltLen:]

	comparisonHash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Constant time comparison
	if len(hash) != len(comparisonHash) {
		return false, nil
	}

	for i := range hash {
		if hash[i] != comparisonHash[i] {
			return false, nil
		}
	}

	return true, nil
}
