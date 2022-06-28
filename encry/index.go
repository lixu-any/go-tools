package lencry

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

//MD5 字符串
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

//SHA1 字符串
func SHA1(text string) string {
	ctx := sha1.New()
	ctx.Write([]byte(text))
	return fmt.Sprintf("%x", ctx.Sum(nil))
}
