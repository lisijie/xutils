package xutils

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

// UniqID 生成一个唯一ID
// 获取一个带前缀、基于当前时间微秒数的唯一ID，类似php的uniqid函数
func UniqID(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%08x%05x", prefix, now.Unix(), now.Nanosecond()/1000)
}

// UUID 生成uuid
func UUID() string {
	return uuid.New().String()
}

// Substr 从字符串s截取指定长度
func Substr(s string, start, length int, padding string) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	if (start + length) > (len(bt) - 1) {
		return string(bt[start:])
	}
	return string(bt[start:start+length]) + padding
}

// StrToInt64 字符串转int64，失败返回0
func StrToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// StrToInt 字符串转int，失败返回0
func StrToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Float64ToStr 返回float64的字符串表示
func Float64ToStr(i float64) string {
	s := strconv.FormatFloat(i, 'E', -1, 64)
	return s
}

// Int64ToStr 返回int64的字符串表示
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

// IntToStr 返回int的字符串表示
func IntToStr(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
