package xcipher

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

func HashPassword(password string) string {

	sha := sha512.New()
	sha.Write([]byte(password))

	code := sha.Sum([]byte{})

	hexHash := strings.ToLower(hex.EncodeToString(code))
	return hexHash
}
