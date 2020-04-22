package service

import (
	"encoding/json"
	"falcon/libs/monitor"
	"github.com/gin-gonic/gin"
	"time"
)

type StopResponse struct {
	Code 		int				`json:"code"`
	Data 		string			`json:"data"`
}


func (s *ServerEntry) HttpStop(ctx *gin.Context) {
	updateTime = monitor.GetUnixNanoTime()-int64(Param.PauseTime)*int64(time.Second)-100

	//设置ContentType
	ctx.Header("Content-Type", "application/json")

	//允许跨域访问
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.JSON(200, GetDemoResponse{0, "ok"})
}

func (h StopResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}
