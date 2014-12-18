package server

import (
	"crypto/md5"
	"encoding/hex"
)

func GetRandomName(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
