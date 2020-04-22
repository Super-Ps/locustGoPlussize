package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type QuitResponse struct {
	Code 		int				`json:"code"`
	Data 		string			`json:"data"`
}

var isQuit bool

func (s *ServerEntry) HttpQuit(ctx *gin.Context) {
	if !isQuit {
		isQuit = true
		deletePidFile()
		go s.quit()
	}

	//设置ContentType
	ctx.Header("Content-Type", "application/json")

	//允许跨域访问
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.JSON(200, GetDemoResponse{0, "ok"})
}

func (h QuitResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}
