package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

type PostDemoResponse struct {
	Code 		int					`json:"code"`
	Data 		*PostDemoData		`json:"data"`
}

type PostDemoData struct {
	Url				string						`json:"url"`
	Header			map[string][]string			`json:"header"`
	ContentType		string						`json:"content_type"`
	ClientIp		string						`json:"client_ip"`
	ByteBody		[]byte						`json:"byte_body"`
	StringBody		string						`json:"string_body"`
	JsonBody		interface{}					`json:"json_body"`
}


func (s *ServerEntry) HttpPostDemo(ctx *gin.Context) {
	//获取url
	url := ctx.Request.URL.String()

	//获取body
	body, _ := ioutil.ReadAll(ctx.Request.Body)

	//获取请求头
	header := ctx.Request.Header

	//获取客户端IP
	clientIP := ctx.ClientIP()

	//获取头的Content-Type(一般用于判断body类型)
	contentType := ctx.ContentType()

	//将body转为json
	var jsonBody interface{}
	if strings.Count(contentType, "application/json") > 0 && len(body) > 0{
		_ = json.Unmarshal([]byte(string(body)), &jsonBody)
	}

	//增加响应头
	ctx.Header("demo", "pc")

	//设置cookies
	ctx.SetCookie("token", "12345678", 10, "/path", "domain", false, false)

	//设置ContentType
	ctx.Header("Content-Type", "application/json")

	//允许跨域访问
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.JSON(200, PostDemoResponse{0, &PostDemoData{
		url,
		header,
		contentType,
		clientIP,
		body,
		string(body),
		jsonBody,
	}})
}


func (h PostDemoResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h PostDemoData) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}