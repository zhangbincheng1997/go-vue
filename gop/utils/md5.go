package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5B ...
func MD5B(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// MD5 ...
func MD5(s string) string {
	return MD5B([]byte(s))
}
