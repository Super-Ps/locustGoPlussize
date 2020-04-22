package service

import (
	"fmt"
	"time"
)

func Printf(msg string, a ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s %s", GetNowTimeString(), msg), a...)
}

func GetNowTimeString() string {
	return time.Unix(0, time.Now().UnixNano()).Format("2006-01-02 15:04:05")
}
