package storage

import (
	"crypto/sha512"
	"encoding/hex"
)

type FilenameConverterFunc func(slug string, key string) string

func EncodeFilenameConverterFunc(slug string, key string) string {
	hash := sha512.Sum512([]byte(slug + key))
	return hex.EncodeToString(hash[:])
}
