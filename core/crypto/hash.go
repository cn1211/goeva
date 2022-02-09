package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

// Deprecated: 直接md5.Sum()会更快,减少一次复制
func MD5Sum(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(data)
	return h.Sum(nil)
}
