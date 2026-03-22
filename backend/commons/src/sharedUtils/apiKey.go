package sharedUtils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

func ComputeAPIKeyHash(rawKey string) (string, error) {
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}
	salt := hex.EncodeToString(saltBytes)
	input := salt + ":" + rawKey
	hash := sha256.Sum256([]byte(input))
	return salt + ":" + hex.EncodeToString(hash[:]), nil
}

func VerifyAPIKeyHash(rawKey string, stored string) bool {
	parts := strings.Split(stored, ":")
	if len(parts) != 2 {
		return false
	}
	salt := parts[0]
	expectedHash := parts[1]
	input := salt + ":" + rawKey
	hash := sha256.Sum256([]byte(input))
	computed := hex.EncodeToString(hash[:])
	return subtle.ConstantTimeCompare([]byte(computed), []byte(expectedHash)) == 1
}
