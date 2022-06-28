package ltime

import (
	"fmt"
	"time"
)

func UninxTime() int64 {
	return time.Now().Unix()
}

//当前日期字符串
func DateNowStr() string {
	tim := time.Unix(UninxTime(), 0).Format("2006-01-02")
	return fmt.Sprintf("%s", tim)
}

//当前日期字符串
func DateTimeNowStr() string {
	tim := time.Unix(UninxTime(), 0).Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s", tim)
}

//当前小时字符串
func DateTimeHour() string {
	tim := time.Unix(UninxTime(), 0).Format("2006::01::02::15")
	return fmt.Sprintf("%s", tim)
}

//当前月
func DateNowMonth() string {
	tim := time.Unix(UninxTime(), 0).Format("2006_01")
	return fmt.Sprintf("%s", tim)
}

//当前年
func DateNowYear() string {
	tim := time.Unix(UninxTime(), 0).Format("2006")
	return fmt.Sprintf("%s", tim)
}

//获取明天的时间戳
func UninxTomTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	return t.Unix() + 1
}

//获取明天的日期
func UninxTomDateStr() string {
	timeStr := time.Now().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	uxi := t.Unix() + 1
	tim := time.Unix(uxi, 0).Format("2006-01-02")
	return tim
}

//获取昨天的日期
func UninxYesDateStr() string {
	timeStr := time.Now().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	tim := time.Unix(t.Unix()-86400, 0).Format("2006-01-02")
	return fmt.Sprintf("%s", tim)
}

//获取今天0点的时间戳
func UninxTodayTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	return t.Unix()
}

//两个日期查数
func TimeDays(startdate, enddate string) float64 {
	a, _ := time.Parse("2006-01-02", startdate)
	b, _ := time.Parse("2006-01-02", enddate)
	d := a.Sub(b)
	return d.Hours() / 24
}

// 计算日期相差多少月
func SubMonth(t11, t22 int64) (month int) {
	t1 := time.Unix(t11, 0)
	t2 := time.Unix(t22, 0)
	y1 := t1.Year()
	y2 := t2.Year()
	m1 := int(t1.Month())
	m2 := int(t2.Month())
	d1 := t1.Day()
	d2 := t2.Day()

	yearInterval := y1 - y2
	// 如果 d1的 月-日 小于 d2的 月-日 那么 yearInterval-- 这样就得到了相差的年数
	if m1 < m2 || m1 == m2 && d1 < d2 {
		yearInterval--
	}
	// 获取月数差值
	monthInterval := (m1 + 12) - m2
	if d1 < d2 {
		monthInterval--
	}
	monthInterval %= 12
	month = yearInterval*12 + monthInterval
	return
}

func GetMonthByUninxTime(t int64) string {
	_time := time.Unix(t, 0)
	month := int(_time.Month())
	monthstr := ""
	if month < 9 {
		monthstr = fmt.Sprintf("0%d", month)
	} else {
		monthstr = fmt.Sprintf("%d", month)
	}
	return fmt.Sprintf("%d_%s", _time.Year(), monthstr)
}

/**
获取本周周一的日期
*/
func GetFirstDateOfWeek() (weekMonday string, uninxtime int64) {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	uninxtime = weekStartDate.Unix()
	weekMonday = weekStartDate.Format("2006-01-02")
	return
}

func GetDaysByUninx(start, end int64) float64 {

	sdate := time.Unix(start, 0).Format("2006-01-02")
	edate := time.Unix(end, 0).Format("2006-01-02")

	return TimeDays(sdate, edate)
}

//当前日期字符串
func GetDateByUninx(utim int64) string {
	tim := time.Unix(utim, 0).Format("2006-01-02")
	return tim
}
