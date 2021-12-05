package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// IsNumber 是否为数字
func IsNumber(s string) bool {
	// 去除首尾空格
	s = strings.TrimSpace(s)
	for i := 0; i < len(s); i++ {
		// 存在 e 或 E, 判断是否为科学计数法
		if s[i] == 'e' || s[i] == 'E' {
			return IsSciNum(s[:i], s[i+1:])
		}
	}
	// 否则判断是否为整数或小数
	return IsInt(s) || IsDec(s)
}

// IsSciNum 是否为科学计数法
func IsSciNum(num1, num2 string) bool {
	// e 前后字符串长度为0 是错误的
	if len(num1) == 0 || len(num2) == 0 {
		return false
	}
	// e 后面必须是整数，前面可以是整数或小数  4  +
	return (IsInt(num1) || IsDec(num1)) && IsInt(num2)
}

// IsDec 判断是否为小数
func IsDec(s string) bool {
	// eg: 11.15, -0.15, +10.15, 3., .15,
	// err: +. 0..
	match1, _ := regexp.MatchString(`^[+-]?\d*\.\d+$`, s)
	match2, _ := regexp.MatchString(`^[+-]?\d+\.\d*$`, s)
	return match1 || match2
}

// IsInt 判断是否为整数
func IsInt(s string) bool {
	match, _ := regexp.MatchString(`^[+-]?\d+$`, s)
	return match
}

type Time struct {
	time.Time
}

const TimeFormat = "2006-1-2 15:4:5"

// MarshalJSON 序列化为JSON
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("\"\""), nil
	}
	stamp := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(stamp), nil
}

// UnmarshalJSON json反序列化为time
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	if len(data) >= 10 {
		t.Time, err = time.ParseInLocation(TimeFormat, FormatTimeString(string(data)), time.Local)
		//t.Time,err = time.Parse(TimeFormat,string(data))
	} else {
		t.Time = time.Time{}
	}
	return err
}

// String 重写String方法
func (t *Time) String() string {
	data, _ := json.Marshal(t)
	return string(data)
}

// SetRaw 读取数据库值
func (t *Time) SetRaw(value interface{}) error {
	switch value.(type) {
	case time.Time:
		t.Time = value.(time.Time)
	}
	return nil
}

// RawValue 写入数据库
func (t *Time) RawValue() interface{} {
	str := t.Format(TimeFormat)
	if str == "0001-01-01 00:00:00" {
		return nil
	}
	return str
}

// FormatTimeString 格式化日期字符串
func FormatTimeString(t string) (ret string) {
	times := strings.ReplaceAll(t, "/", "-")
	arr := strings.Split(times, " ")
	if len(arr) == 1 || len(arr) == 0 {
		ret = strings.Join([]string{arr[0], "00:00:00"}, " ")
	} else {
		switch strings.Count(arr[1], ":") {
		case 0:
			ret = strings.Join([]string{arr[0], strings.Join([]string{arr[1], ":00:00"}, "")}, " ")
			break
		case 1:
			ret = strings.Join([]string{arr[0], strings.Join([]string{arr[1], ":00"}, "")}, " ")
			break
		default:
			ret = times
			break
		}
	}
	return
}

// ConvertToFormatDay 转换日期
func ConvertToFormatDay(excelDaysString string) string {
	// 2006-01-02 距离 1900-01-01的天数
	baseDiffDay := 38719 //在网上工具计算的天数需要加2天，什么原因没弄清楚
	curDiffDay := excelDaysString
	b, _ := strconv.Atoi(curDiffDay)
	// 获取excel的日期距离2006-01-02的天数
	realDiffDay := b - baseDiffDay
	//fmt.Println("realDiffDay:",realDiffDay)
	// 距离2006-01-02 秒数
	realDiffSecond := realDiffDay * 24 * 3600
	//fmt.Println("realDiffSecond:",realDiffSecond)
	// 2006-01-02 15:04:05距离1970-01-01 08:00:00的秒数 网上工具可查出
	baseOriginSecond := 1136185445
	resultTime := time.Unix(int64(baseOriginSecond+realDiffSecond), 0).Format("2006-01-02")
	return resultTime
}

// MemberInArray 元素是否在数组内
func MemberInArray(target string, array []string) bool {
	sort.Strings(array)
	index := sort.SearchStrings(array, target)
	if index < len(array) && array[index] == target {
		return true
	}
	return false
}

// ArrayContainMember 数组是否包含元素
func ArrayContainMember(target string, array []string) bool {
	for _, ele := range array {
		if strings.Contains(target, ele) {
			return true
		}
	}
	return false
}
