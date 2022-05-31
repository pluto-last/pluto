package utils

import "time"

func GetDayRange(t time.Time) (start, end time.Time) {
	loc, _ := time.LoadLocation("Local")
	date := t.Format("2006-01-02")
	start, _ = time.ParseInLocation("2006-01-02 15:04:05", date+" 00:00:00", loc)
	end, _ = time.ParseInLocation("2006-01-02 15:04:05", date+" 23:59:59", loc)
	return
}

// GetLoginTimeOut 获取登陆超时时间，到当天0点
func GetLoginTimeOut() int64 {
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)
	remainSecond := time.Duration(todayLastTime.Unix()-time.Now().Local().Unix()) * time.Second
	minutes := int64(remainSecond.Minutes())
	return minutes
}
