package lconv

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/astaxie/beego/logs"
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

//字符串转换位int
func ParseInt64(val interface{}) int64 {

	if val == nil {
		return 0
	}

	switch ret := val.(type) {
	case string:
		return StrToInt64(ret)
	case int8:
		return int64(ret)
	case uint8:
		return int64(ret)
	case int16:
		return int64(ret)
	case uint16:
		return int64(ret)
	case int:
		return int64(ret)
	case uint:
		return int64(ret)
	case int64:
		return int64(ret)
	case uint64:
		return int64(ret)
	case float32:
		return int64(ret)
	case float64:
		return int64(ret)
	}

	return 0
}

func ParseInt(val interface{}) int {

	if val == nil {
		return 0
	}

	switch ret := val.(type) {
	case string:
		return StrToInt(ret)
	case int8:
		return int(ret)
	case uint8:
		return int(ret)
	case int16:
		return int(ret)
	case uint16:
		return int(ret)
	case int:
		return int(ret)
	case uint:
		return int(ret)
	case int64:
		return int(ret)
	case uint64:
		return int(ret)
	case float32:
		return int(ret)
	case float64:
		return int(ret)
	}

	return 0
}

//浮点数转换int卡类型
func FloatTostr(floatstr float64, i int) string {
	return strconv.FormatFloat(floatstr, 'E', -1, i)
}

//抓换成json
func EveryToJson(data interface{}) map[string]interface{} {
	str, _ := json.Marshal(data)
	newval := make(map[string]interface{})
	err := json.Unmarshal(str, &newval)
	if err != nil {
		logs.Error("EveryToJson:", err)
	}
	return newval
}

//任意类型转换字符串
func InterfaceToStr(v interface{}) string {
	str := ""
	switch v.(type) {
	case string:
		str = v.(string)
	default:
		jsonbyte, _ := json.Marshal(v)
		str = string(jsonbyte)
	}
	return str
}

//截取字符串
func Substr(str string, l int) string {
	if len(str) > l {
		str2 := []rune(str)
		return string(str2[0:l])
	}
	return str
}

//json转字符串
func JsonEncode(data interface{}) string {
	datastr, _ := json.Marshal(data)
	return string(datastr)
}

func StructMap(obj interface{}) map[string]interface{} {

	var data = make(map[string]interface{})

	by, _ := json.Marshal(obj)

	json.Unmarshal([]byte(by), &data)

	return data
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func ToInt(a interface{}) int {

	ai := 0

	switch a := a.(type) {
	case float64:
		ai = int(a)
	case float32:
		ai = int(a)
	case int:
		ai = a
	case int32:
		ai = int(a)
	case int64:
		ai = int(a)
	case string:
		ai = StrToInt(a)

	default:

	}

	return ai
}

func ToInt64(a interface{}) int64 {

	var ai int64

	switch a := a.(type) {
	case float64:
		ai = int64(a)
	case float32:
		ai = int64(a)
	case int:
		ai = int64(a)
	case int32:
		ai = int64(a)
	case int64:
		ai = a
	case string:
		ai = StrToInt64(a)

	default:

	}

	return ai
}

func ToFloat64(a interface{}) float64 {

	var ai float64

	switch a := a.(type) {
	case float64:
		ai = a
	case float32:
		ai = float64(a)
	case int:
		ai = float64(a)
	case int32:
		ai = float64(a)
	case int64:
		ai = float64(a)
	case string:
		ai = StrToFloat64(a, 64)

	default:

	}

	return ai
}
