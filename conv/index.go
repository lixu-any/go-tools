package lconv

import (
	"encoding/json"
	"strconv"
)

//数字转换字符串
func IntToStr(i int) string {
	return strconv.Itoa(i)
}

//64位数字转换字符串
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

//字符串转换位int
func StrToInt(str string) int {
	intnum, _ := strconv.Atoi(str)
	return intnum
}

//字符串转换位int64
func StrToInt64(str string) int64 {
	intnum, _ := strconv.ParseInt(str, 10, 64)
	return intnum
}

//字符串转换浮点数
func StrToFloat64(str string, i int) float64 {
	v2, _ := strconv.ParseFloat(str, i)
	return v2
}

//浮点数转换int卡类型
func FloatTostr(floatstr float64, i int) string {
	return strconv.FormatFloat(floatstr, 'E', -1, i)
}

//json转字符串
func JsonEncode(data interface{}) string {
	datastr, _ := json.Marshal(data)
	return string(datastr)
}

//截取字符串
func Substr(str string, l int) string {
	if len(str) > l {
		str2 := []rune(str)
		return string(str2[0:l])
	}
	return str
}
