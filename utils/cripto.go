package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// GetHashPassword ...
func GetHashPassword(s string) (string, error) {
	data := []byte(s)
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return "", err
	}
	s = hex.EncodeToString(hash.Sum(nil))
	return s, nil
}
