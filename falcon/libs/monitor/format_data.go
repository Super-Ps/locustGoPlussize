package monitor

import "time"

func MilliNow() int64 {
	return GetUnixNanoTime()/int64(time.Millisecond)
}

func TakingMilli(startTime int64) int64 {
	return (MilliNow()-startTime)
}

func FormatAbsInt64(s int64) int64 {
	if s < 0 {
		return s*int64(-1)
	} else {
		return s
	}
}

func FormatNegativeInt64(s int64) int64 {
	if s < 0 {
		return int64(0)
	} else {
		return s
	}
}
