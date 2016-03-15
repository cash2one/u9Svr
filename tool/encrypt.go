package tool

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func HmacSHA1Encrypt(context, secKey string) []byte {
	mac := hmac.New(sha1.New, []byte(secKey))
	mac.Write([]byte(context))
	return mac.Sum(nil)
}
