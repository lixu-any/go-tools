package larray

//判断字符串在数组内
func InArray(s string, d []string) bool {
	for _, v := range d {
		if s == v {
			return true
		}
	}
	return false
}
