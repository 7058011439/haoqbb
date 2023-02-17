package Timer

import "time"

func GetCurrDateString(format string) string {
	if format == "" {
		format = "2006-01-02"
	}
	timestamp := time.Now().Unix()

	tm := time.Unix(timestamp, 0)

	return tm.Format(format)
}

// 获取系统时间，精确到毫秒
func GetOsTimeMillisecond() int64 {
	return time.Now().Local().UnixNano() / int64(time.Millisecond)
}

// 获取系统时间，精确到秒
func GetOsTimeSecond() int32 {
	return int32(time.Now().Local().Unix())
}
