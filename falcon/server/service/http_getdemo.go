package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type GetDemoResponse struct {
	Code 		int				`json:"code"`
	Data 		string			`json:"data"`
}


func (s *ServerEntry) HttpGetDemo(ctx *gin.Context) {
	//设置ContentType
	ctx.Header("Content-Type", "application/json")

	//允许跨域访问
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.JSON(200, GetDemoResponse{0, "ok"})
}

func (h GetDemoResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}