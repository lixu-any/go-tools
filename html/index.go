package l_html

import (
	"regexp"
	"strings"
)

func NOHTML(content interface{}) string {
	str := content.(string)
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllStringFunc(str, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	str = re.ReplaceAllString(str, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	str = re.ReplaceAllString(str, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllString(str, "")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	str = re.ReplaceAllString(str, "")

	str = strings.TrimSpace(str)

	str = strings.Replace(str, "\"", "", -1)
	str = strings.Replace(str, "&nbsp;", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "'", "", -1)

	return str
}
