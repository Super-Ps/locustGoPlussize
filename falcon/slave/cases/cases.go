package cases

import (
	"falcon/slave/site/boomer"
	"falcon/slave/site/boomer_client"
	"encoding/json"
)


//入口配置（不能改名），必须包含*boomer.WorkerEntry，其他字段可选（可自行新增需要的字段）
type CaseEntry struct {
	*boomer.WorkerEntry
	http 				*boomer_client.HttpWorker
	fastHttp 			*boomer_client.FastHttpWorker
	httpGetConfig		*boomer_client.HttpRequestConfig
	httpPostConfig		*boomer_client.HttpRequestConfig
	fastHttpGetConfig	*boomer_client.FastHttpRequestConfig
	fastHttpPostConfig	*boomer_client.FastHttpRequestConfig
	httpProxy 			*boomer_client.HttpWorker			//http代理client
	fastHttpProxy 		*boomer_client.FastHttpWorker		//http代理client
}

//初始化函数（不能改名），可自定义内容，但必须返回&CaseEntry{
func NewCase() interface{} {
	http := boomer_client.NewHttpClient()
	fastHttp := boomer_client.NewFastHttpClient()
	httpProxy := boomer_client.NewDefaultProxyHttpClient(boomer.Conf.BaseConfig.Custom.HttpProxy)
	fastHttpProxy := boomer_client.NewDefaultProxyFastHttpClient(boomer.Conf.BaseConfig.Custom.HttpProxy)
	return &CaseEntry{
		http: http,
		fastHttp: fastHttp,
		httpProxy: httpProxy,
		fastHttpProxy: fastHttpProxy,
	}
}

//入口函数（不能改名），worker启动时，运行一次
func (c *CaseEntry) OnStart() {
	url := c.Config.Custom.Url
	getRoute := c.Config.Custom.GetRoute
	postRoute := c.Config.Custom.PostRoute
	contentType := c.Config.Custom.ContentType
	getParams := nilMap()
	postParams := map[string]string{
		"test": "abc",
	}
	getHeaders := nilMap()

	postHeaders := map[string]string{
		"Cookie": "token=123",
		"Wps-Sid": "[{\"key\":\"Wps-Sid\",\"value\":\"0987654321\",\"description\":\"\",\"type\":\"text\",\"enabled\":true}]",
	}
	getBody := nilBody()
	postBody, _ := json.Marshal(struct{
		Test1 	string  	`json:"test1"`
	}{Test1: "abcd1"})

	c.httpGetConfig = &boomer_client.HttpRequestConfig{
		Url: url,
		Route: getRoute,
		ContentType: contentType,
		Params: getParams,
		Headers: getHeaders,
		Body: getBody,
		Request: nil,
	}

	c.httpPostConfig = &boomer_client.HttpRequestConfig{
		Url: url,
		Route: postRoute,
		ContentType: contentType,
		Params: postParams,
		Headers: postHeaders,
		Body: postBody,
		Request: nil,
	}

	c.fastHttpGetConfig = &boomer_client.FastHttpRequestConfig{
		Url: url,
		Route: getRoute,
		ContentType: contentType,
		Params: getParams,
		Headers: getHeaders,
		Body: getBody,
		Request: nil,
	}

	c.fastHttpPostConfig = &boomer_client.FastHttpRequestConfig{
		Url: url,
		Route: postRoute,
		ContentType: contentType,
		Params: postParams,
		Headers: postHeaders,
		Body: postBody,
		Request: nil,
	}
}

//退出函数（不能改名），worker停止前，运行一次
func (c *CaseEntry) OnStop() {
	index := c.Index		//Worker在本Slave的序号
	num := c.Num			//Worker在所有Slave中的序号
	times := c.Times		//Worker在本Slave的当前运行次数
	c.Log.Info("Worker即将退出! Worker在本Slave(%s-%s)的序号:%d, Worker在所有Slave中的序号:%d, Worker在本Slave的已运行次数:%d", c.GetMachineID(), c.GetSlaveID(), index, num, times)
}

//空map函数
func nilMap() map[string]string {
	return nil
}

//空Body函数
func nilBody() []byte {
	return nil
}
