package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)


type AddressResponse struct {
	Code 		int				`json:"code"`
	Data 		*AddressData	`json:"data"`
}

type AddressData struct {
	Ip 			string 			`json:"ip"`
}

func (s *ServerEntry) HttpAddress(ctx *gin.Context) {
	//设置ContentType
	ctx.Header("Content-Type", "application/json")

	//允许跨域访问
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.JSON(200, AddressResponse{0, &AddressData{GetClientIP(ctx)}})
}

func GetClientIP(ctx *gin.Context) string {
	var forwardIp string
	clientIp := strings.TrimSpace(ctx.ClientIP())
	forwardIpSplit := strings.Split(strings.TrimSpace(ctx.Request.Header.Get("X-Forwarded-For")), ",")
	realIp := strings.TrimSpace(ctx.Request.Header.Get("X-real-ip"))
	for _, v := range forwardIpSplit {
		if v == "" || v == "127.0.0.1" {
			continue
		}
		forwardIp = v
	}
	if forwardIp != "" {
		return forwardIp
	} else if realIp != "" && realIp != "127.0.0.1" {
		return realIp
	} else if clientIp != "" && clientIp != "127.0.0.1" {
		return clientIp
	} else {
		return ""
	}
}

func (h AddressResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h AddressData) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}