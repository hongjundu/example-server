package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Encode(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))

	cipherStr := hash.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
