package hasher

import (
	"crypto/md5"
	"encoding/hex"
)

func GeneratePasswordHash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
