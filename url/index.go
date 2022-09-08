package lurl

import (
	"crypto/tls"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/astaxie/beego/httplib"
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

func PostBody(apiurl string, data map[string]interface{}) (content []byte, httpcode int, err error) {

	req := httplib.Post(apiurl)

	if strings.Contains(apiurl, "https://") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	req.JSONBody(data)

	resp, err := req.Response()

	if err != nil {
		return
	}

	httpcode = resp.StatusCode

	content, err = ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	return
}

func PostParam(apiurl string, data map[string]string) (content []byte, httpcode int, err error) {

	req := httplib.Post(apiurl)

	if strings.Contains(apiurl, "https://") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	for k, v := range data {
		req.Param(k, v)
	}

	resp, err := req.Response()

	if err != nil {
		return
	}

	httpcode = resp.StatusCode

	content, err = ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	return
}

func Get(apiurl string) (content []byte, httpcode int, err error) {

	req := httplib.Get(apiurl)

	if strings.Contains(apiurl, "https://") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	resp, err := req.Response()

	if err != nil {
		return
	}

	httpcode = resp.StatusCode

	content, err = ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	return
}
