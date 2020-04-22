

package cases

import (
"encoding/json"
"errors"
"falcon/slave/site/boomer_client"
"fmt"
)

//定义Post send body结构体
type PostBody struct {
	Password	string		`json:"password"`
}

//定义Get send body结构体
type GetBody struct {
	Password	string		`json:"password"`
}

//定义CustomPost函数返回结构体
type PostFuncRes struct {
	ClientType	string
	Ip 			string
}

//定义请求的api的restful格式
type RestFulBody struct {
	Code 		int						`json:"code"`
	Data		DataBody				`json:"data"`

}

type DataBody struct {
	Url 			string 					`json:"url"`
	Header			map[string][]string		`json:"header"`
	ContentType		string					`json:"content_type"`
	ClientIp		string					`json:"client_ip"`
	ByteBody		string					`json:"byte_body"`
	StringBody		string					`json:"string_body"`
	JsonBody		interface{}				`json:"json_body"`
}


func (c *CaseEntry) CustomDemo() {
	postRes, err := c.CustomPost()
	if err != nil {
		return
	}
	c.CustomGet(postRes)
}

func (c *CaseEntry) CustomPost() (*PostFuncRes, error) {
	//设置参数
	userId := fmt.Sprintf("user_%d", c.Num)
	pwd := fmt.Sprintf("pwd_%d", c.Num)
	name := fmt.Sprintf("name_%d", c.Num)
	header := map[string]string{
		"Cookie": userId,
	}
	params := map[string]string{
		"name": name,
	}
	body := PostBody{
		Password: pwd,
	}
	sendBody, _ := json.Marshal(body)

	//初始化请求配置
	resConf := &boomer_client.HttpRequestConfig{
		Url: c.Config.Custom.Url,
		Route: c.Config.Custom.PostRoute,
		ContentType: c.Config.Custom.ContentType,
		Params: params,
		Headers: header,
		Body: sendBody,
		Request: nil,
	}

	//发送请求, isCatch=true，则需要自行判断返回结果并释放请求
	resp, err := c.http.Post(resConf, true)
	if err != nil {
		c.Log.Error(fmt.Sprintf("Post %s fail:%s", c.Config.Custom.PostRoute, err.Error()))
		c.RecordFailure("POST", c.Config.Custom.PostRoute, 0, err.Error())
		return nil, err
	}
	//函数return前自动释放请求
	defer c.http.FreeHttpResult(resp)

	//判断返回码
	if resp.Code != 200 {
		c.RecordFailure("POST", c.Config.Custom.PostRoute, int64(resp.UsedMillTimes), fmt.Sprintf("Response code: %d", resp.Code))
		return nil, errors.New(fmt.Sprintf("Response code: %d", resp.Code))
	}
	c.RecordSuccess("POST", c.Config.Custom.PostRoute, int64(resp.UsedMillTimes), resp.Length)

	//获取header
	respHeader := resp.Header
	clientType := respHeader.Get("Demo")
	if clientType == "" {
		c.RecordFailure("POST", c.Config.Custom.PostRoute, int64(resp.UsedMillTimes), fmt.Sprintf("Get header fail"))
		return nil, errors.New(fmt.Sprintf("Get header fail"))
	}

	//获取body
	var jsonBody RestFulBody
	respBody := resp.Body
	err = json.Unmarshal(respBody, &jsonBody)
	if err != nil {
		c.RecordFailure("POST", c.Config.Custom.PostRoute, int64(resp.UsedMillTimes), fmt.Sprintf("Unmarshal fail: %s", err.Error()))
		return nil, errors.New(fmt.Sprintf("Unmarshal fail: %s", err.Error()))
	}
	ip := jsonBody.Data.ClientIp

	return &PostFuncRes{
		ClientType: clientType,
		Ip: ip,
	}, nil
}

func (c *CaseEntry) CustomGet(postRes *PostFuncRes) {
	ip := postRes.Ip
	clientType := postRes.ClientType
	params := map[string]string{
		"ip": ip,
	}
	header := map[string]string{
		"type": clientType,
	}
	resConf := &boomer_client.HttpRequestConfig{
		Url: c.Config.Custom.Url,
		Route: c.Config.Custom.GetRoute,
		ContentType: c.Config.Custom.ContentType,
		Params: params,
		Headers: header,
		Body: nil,
		Request: nil,
	}
	_, _ = c.http.Get(resConf)
}
