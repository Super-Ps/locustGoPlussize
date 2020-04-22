package cases

import (
	"encoding/json"
	"falcon/slave/site/boomer_client"
	"fmt"
	"strings"
)

func (c *CaseEntry) FastHttpGetDemo() {
	_, _ = c.fastHttp.Get(c.fastHttpGetConfig)
}

func (c *CaseEntry) FastHttpPostDemo() {
	resp, err := c.fastHttp.Post(c.fastHttpPostConfig, true)
	c.fastHttpPost(resp, err, fmt.Sprintf("FastHttp(%s)", c.httpPostConfig.Route))
}

func (c *CaseEntry) FastHttpProxyDemo() {
	resp, err := c.fastHttpProxy.Post(c.fastHttpPostConfig, true)
	c.fastHttpPost(resp, err, fmt.Sprintf("FastHttpProxy(%s)", c.httpPostConfig.Route))
}

func (c *CaseEntry) fastHttpPost(resp *boomer_client.FastHttpWorkerResponse, respErr error, name string) {
	if respErr != nil {
		c.Log.Error(fmt.Sprintf("Post %s fail:%s", c.httpPostConfig.Route, respErr.Error()))
		c.RecordFailure("POST", name, int64(resp.UsedMillTimes), respErr.Error())
		return
	}
	defer c.fastHttp.FreeFastHttpResult(resp)

	if resp.Code != 200 {
		c.RecordFailure("POST", name, int64(resp.UsedMillTimes), fmt.Sprintf("Response code: %d", resp.Code))
		return
	}
	c.RecordSuccess("POST", name, int64(resp.UsedMillTimes), resp.Length)

	header := resp.Header
	headerWpsClient := string(header.Peek("Wps-Client"))

	var jsonBody map[string]interface{}
	body := resp.Body
	_ = json.NewDecoder(strings.NewReader(string(body))).Decode(&jsonBody)
	bodyCode := jsonBody["code"]
	bodyContent := jsonBody["content"]

	logData := fmt.Sprintf("%s - httpcode:%+v\nheader:%+v\nWpsClient:%+v\nbody:%+v\n,bodyCode:%+v\nbodyContent:%+v\n", name, resp.Code, header, headerWpsClient, jsonBody, bodyCode, bodyContent)
	_ = logData
	//c.Log.Info(logData)
}
