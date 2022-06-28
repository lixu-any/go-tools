package lurl

import (
	"net/url"
	"strings"
)

func UrlEncode(str string) string {
	encodeurl := url.QueryEscape(str)
	encodeurl = strings.ReplaceAll(encodeurl, "\"", "%22")
	return encodeurl
}

func UrlDecode(str string) string {
	decodeurl, _ := url.QueryUnescape(str)
	return decodeurl
}
