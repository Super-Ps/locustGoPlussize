package service

import (
	"encoding/json"
	"falcon/libs/monitor"
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
)

var exePath, _ = monitor.GetProcessExePath(os.Getpid())

type RestartResponse struct {
	Code 		int				`json:"code"`
	Data 		string			`json:"data"`
}

var isRestart bool

func (s *ServerEntry) HttpRestart(ctx *gin.Context) {
	if !isRestart {
		isRestart = true
		deletePidFile()
		go s.restart()
	}

	//设置ContentType
	ctx.Header("Content-Type", "application/json")

	//允许跨域访问
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.JSON(200, GetDemoResponse{0, "ok"})
}

func (h RestartResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (s *ServerEntry) restart() {
	dir, _ := os.Getwd()
	cmd := &exec.Cmd{
		Path: exePath,
		Args: os.Args,
		Dir: dir,
	}
	_ = cmd.Start()
	go s.quit()
}
