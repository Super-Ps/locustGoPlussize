package cases

import (
	"encoding/json"
	"falcon/slave/site/boomer_client"
	"fmt"
	"strings"
)

func (c *CaseEntry) HttpGetDemo() {

	res, err := c.http.Get(c.httpGetConfig, true)
	if err == nil {
		defer c.http.FreeHttpResult(res)
		c.RecordSuccess("GET", c.httpGetConfig.Route, int64(res.UsedMillTimes), res.Length)
		c.Log.Info("status:%s, header:%s, body:%s", res.Code, res.Header, res.Body)
		fmt.Printf("status:%s, header:%s, body:%s", res.Code, res.Header, res.Body)
	} else {
		c.RecordFailure("GET", c.httpGetConfig.Route, 0, err.Error())
		c.Log.Error("error:%s", err)
	}
	//HttpWorkerResponse, error := c.http.Get(c.httpGetConfig)
	//fmt.Printf("????GET请求返回值",HttpWorkerResponse,error,"\n????类型是,",reflect.TypeOf(HttpWorkerResponse))
}

func (c *CaseEntry) HttpPostDemo() {
	resp, err := c.http.Post(c.httpPostConfig, true)
	c.httpPost(resp, err, fmt.Sprintf("Http(%s)", c.httpPostConfig.Route))
}

func (c *CaseEntry) HttpProxyDemo() {
	resp, err := c.httpProxy.Post(c.httpPostConfig, true)
	c.httpPost(resp, err, fmt.Sprintf("HttpProxy(%s)", c.httpPostConfig.Route))
}

func (c *CaseEntry) httpPost(resp *boomer_client.HttpWorkerResponse, respErr error, name string) {
	if respErr != nil {
		c.Log.Error(fmt.Sprintf("Post %s fail:%s", c.httpPostConfig.Route, respErr.Error()))
		c.RecordFailure("POST", name, int64(resp.UsedMillTimes), respErr.Error())
		return
	}
	defer c.http.FreeHttpResult(resp)

	if resp.Code != 200 {
		c.RecordFailure("POST", name, int64(resp.UsedMillTimes), fmt.Sprintf("Response code: %d", resp.Code))
		return
	}
	c.RecordSuccess("POST", name, int64(resp.UsedMillTimes), resp.Length)

	header := resp.Header
	headerWpsClient := header.Get("Wps-Client")

	var jsonBody map[string]interface{}
	body := resp.Body
	_ = json.NewDecoder(strings.NewReader(string(body))).Decode(&jsonBody)
	bodyCode := jsonBody["code"]
	bodyContent := jsonBody["content"]

	logData := fmt.Sprintf("%s - httpcode:%+v\nheader:%+v\nWpsClient:%+v\nbody:%+v\n,bodyCode:%+v\nbodyContent:%+v\n", name, resp.Code, header, headerWpsClient, jsonBody, bodyCode, bodyContent)
	_ = logData
	//c.Log.Info(logData)
}

