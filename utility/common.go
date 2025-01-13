package utility

import "time"

// UnixMilli dùng để lấy mili giây của thời gian cho trước
// @params - thời gian
// @returns - mili giây của thời gian cho trước
func UnixMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// CurrentTimeInMilli dùng để lấy thời gian hiện tại tính bằng mili giây
// Hàm này sẽ được sử dụng khi cần timestamp hiện tại
// @returns - timestamp hiện tại (tính bằng mili giây)
func CurrentTimeInMilli() int64 {
	return UnixMilli(time.Now())
}
