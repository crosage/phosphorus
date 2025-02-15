package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

var salt = "rtrt"

func GeneratePassHash(password string) string {
	password += salt
	hasher := sha256.New()
	hasher.Write([]byte(password))
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha
}
