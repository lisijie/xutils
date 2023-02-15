package xutils

import (
	"strings"
	"time"
)

var datePatterns = []string{
	// year
	"Y", "2006",
	"y", "06",

	// month
	"m", "01",
	"n", "1",
	"M", "Jan",
	"F", "January",

	// day
	"d", "02",
	"j", "2",

	// week
	"D", "Mon",
	"l", "Monday",

	// time
	"g", "3",
	"G", "15",
	"h", "03",
	"H", "15",

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// DateParse 解析时间字符串
func DateParse(dateString, format string, locs ...*time.Location) (time.Time, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	if len(locs) > 0 {
		return time.ParseInLocation(format, dateString, locs[0])
	} else {
		return time.ParseInLocation(format, dateString, time.Local)
	}
}

// Date 格式化时间
func Date(t time.Time, format string) string {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return t.Format(format)
}

// UnixDate 将unix时间戳格式化为时间格式
func UnixDate(ts int64, format string) string {
	return Date(time.Unix(ts, 0), format)
}

// DateUnix 解析时间字符串为unix时间戳
func DateUnix(dateString, format string) int64 {
	t, err := DateParse(dateString, format)
	if err != nil {
		return 0
	} else {
		return t.Unix()
	}
}
