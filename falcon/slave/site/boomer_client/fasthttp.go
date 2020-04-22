package boomer_client

import (
	"bufio"
	"crypto/md5"
	"errors"
	"falcon/slave/site/boomer"
	"fmt"
	"github.com/valyala/fasthttp"
	"net"
	"strconv"
	"strings"
	"time"
)

//请求配置
type FastHttpRequestConfig struct {
	Url     		string
	Route   		string
	ContentType		string
	Method 			string
	Params  		map[string]string
	Headers 		map[string]string
	Body    		[]byte
	Request 		*fasthttp.Request
	Blocks			int
	BlocksWait		int64
	NoRespBody		bool
	Close 			bool
	GetMd5			bool
	WorkerResponse 	func(ctx *FastHttpWorkerResponse)
}

//返回值结构
type FastHttpWorkerResponse struct {
	Request        *fasthttp.Request
	Response       *fasthttp.Response
	BlockResult		*FastHttpResponseResult
	*FastHttpResponseResult
}

type FastHttpResponseResult struct {
	Body   			[]byte
	Code   			int
	Length 			int64
	Size 			int64
	Header 			fasthttp.ResponseHeader
	Md5				string
	ByteSpeed		float64
	Error 			error
	TotalMillTimes  float64
	UsedMillTimes  	float64
	StartNanoTime	int64
	EndNanoTime		int64
	UsedNanoTime	int64
	TotalNanoTime	int64
}

type FastHttpWorker struct {
	Client *FastHttpWorkerClient
}

type FastHttpWorkerClient struct {
	*fasthttp.Client
}

var DefaultFastHttpBlocks = 1024*1024


//创建Client
func NewFastHttpClient(c ... *FastHttpWorkerClient) *FastHttpWorker {
	defaultClient := &FastHttpWorkerClient{
		Client: &fasthttp.Client{
			ReadTimeout:               time.Second * 60,
			WriteTimeout:              time.Second * 60,
			MaxIdleConnDuration:       time.Second * 60,
			NoDefaultUserAgentHeader:  true,
			MaxConnsPerHost:           50000,
			MaxIdemponentCallAttempts: 50000,
		},
	}

	var client FastHttpWorkerClient
	if len(c) == 0 {
		client = *defaultClient
	} else {
		client = *c[0]
	}
	return &FastHttpWorker{Client: &client}
}

//创建默认带Http代理的Client
func NewDefaultProxyFastHttpClient(proxyAddr string) *FastHttpWorker {
	defaultClient := &FastHttpWorkerClient{
		Client: &fasthttp.Client{
			Dial: FastHttpDialer(proxyAddr),
			ReadTimeout:               time.Second * 60,
			WriteTimeout:              time.Second * 60,
			MaxIdleConnDuration:       time.Second * 60,
			NoDefaultUserAgentHeader:  true,
			MaxConnsPerHost:           50000,
			MaxIdemponentCallAttempts: 50000,
		},
	}
	return &FastHttpWorker{Client: defaultClient}
}

//Http代理函数
func FastHttpDialer(proxyAddr string) fasthttp.DialFunc {
	return func(addr string) (net.Conn, error) {
		conn, err := fasthttp.Dial(proxyAddr)
		if err != nil {
			return nil, err
		}

		req := "CONNECT " + addr + " HTTP/1.1\r\n"

		req += "\r\n"

		if _, err := conn.Write([]byte(req)); err != nil {
			return nil, err
		}

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		res.SkipBody = true

		if err := res.Read(bufio.NewReader(conn)); err != nil {
			_ = conn.Close()
			return nil, err
		}
		if res.Header.StatusCode() != 200 {
			_ = conn.Close()
			return nil, fmt.Errorf("could not connect to proxy")
		}
		return conn, nil
	}
}

//发送Post请求, isCatch=False,自动向Master发送事件
func (w *FastHttpWorker) Post(res *FastHttpRequestConfig, isCatch ...bool) (*FastHttpWorkerResponse, error) {
	return w.Do("POST", res, isCatch...)
}

//发送Get请求
func (w *FastHttpWorker) Get(res *FastHttpRequestConfig, isCatch ...bool) (*FastHttpWorkerResponse, error) {
	return w.Do("GET", res, isCatch...)
}

//发送Put请求
func (w *FastHttpWorker) Put(res *FastHttpRequestConfig, isCatch ...bool) (*FastHttpWorkerResponse, error) {
	return w.Do("PUT", res, isCatch...)
}

//发送Delete请求
func (w *FastHttpWorker) Delete(res *FastHttpRequestConfig, isCatch ...bool) (*FastHttpWorkerResponse, error) {
	return w.Do("DELETE", res, isCatch...)
}

//发送请求
func (w *FastHttpWorker) sendFastHttpRequest(method string, resConf *FastHttpRequestConfig) (*FastHttpWorkerResponse, error) {
	result := &FastHttpWorkerResponse{
		Request: nil,
		Response: nil,
		BlockResult: &FastHttpResponseResult{},
		FastHttpResponseResult: &FastHttpResponseResult{},
	}

	url := strings.Trim(resConf.Url, "/") + resConf.Route
	split := "?"
	for k, v := range resConf.Params {
		url += fmt.Sprintf("%s%s=%s", split, k, v)
		split = "&"
	}

	if resConf.Request == nil {
		result.Request = fasthttp.AcquireRequest()
		result.Request.Header.SetMethod(method)
		result.Request.SetRequestURI(url)

		for k, v := range resConf.Headers {
			result.Request.Header.Set(k, v)
		}

		if resConf.ContentType != "" {
			result.Request.Header.SetContentType(resConf.ContentType)
		}

		if resConf.Body != nil {
			result.Request.SetBody(resConf.Body)
		}
	} else {
		result.Request = resConf.Request
	}
	if resConf.Blocks <= 0 {
		resConf.Blocks = DefaultHttpBlocks
	}
	if resConf.WorkerResponse == nil {
		resConf.WorkerResponse = DefaultFastWorkerResponseFunc
	}
	if resConf.Close {
		result.Request.SetConnectionClose()
		result.Request.ConnectionClose()
	}
	result.StartNanoTime = time.Now().UnixNano()
	result.BlockResult.StartNanoTime = result.StartNanoTime


	var doErr error
	result.Response = fasthttp.AcquireResponse()
	doErr = w.Client.Do(result.Request, result.Response)
	result.EndNanoTime = time.Now().UnixNano()
	result.BlockResult.EndNanoTime = result.EndNanoTime
	result.UsedNanoTime += (time.Now().UnixNano() - result.StartNanoTime)
	if doErr != nil {
		return result, doErr
	}

	result.BlockResult.StartNanoTime = time.Now().UnixNano()
	buf := result.Response.Body()
	result.BlockResult.UsedNanoTime = time.Now().UnixNano() - result.BlockResult.StartNanoTime
	result.UsedNanoTime += result.BlockResult.UsedNanoTime
	if len(buf) > 0 {
		if resConf.GetMd5 {
			md5AllBuf := md5.New()
			md5AllBuf.Write(buf)
			result.BlockResult.Md5 = strings.ToUpper(fmt.Sprintf("%x",md5AllBuf.Sum(nil)))
			result.Md5 = result.BlockResult.Md5
		}
		if !resConf.NoRespBody {
			result.BlockResult.Body = buf
			result.Body = buf
		}
	}
	result.BlockResult.Size = int64(len(buf))
	result.Size = result.BlockResult.Size
	w.UpdateFastHttpResult(result)
	resConf.WorkerResponse(result)
	time.Sleep(time.Duration(resConf.BlocksWait)*time.Millisecond)
	return result, nil
}

//创建请求
func (w *FastHttpWorker) Do(method string, res *FastHttpRequestConfig, isCatch ...bool) (*FastHttpWorkerResponse, error) {
	boom := boomer.Slave.Boom
	param := boomer.Param
	minWait := param.MinWait
	maxWait := param.MaxWait
	defer boomer.RandomWait(time.Duration(minWait)*time.Millisecond, time.Duration(maxWait)*time.Millisecond)

	var result *FastHttpWorkerResponse
	var err error

	result, err = w.sendFastHttpRequest(method, res)

	if len(isCatch) > 0 {
		if isCatch[0] {
			return result, err
		}
	}

	defer w.FreeFastHttpResult(result)
	if err != nil {
		boom.RecordFailure(method, res.Route, 0, err.Error())
		return result, err
	}

	if (result.Code < 200) || (result.Code >= 300) {
		err := errors.New(fmt.Sprintf("Response code is %d", result.Code))
		boom.RecordFailure(method, res.Route, int64(result.UsedMillTimes), err.Error())
		return result, err
	}
	boom.RecordSuccess(method, res.Route, int64(result.UsedMillTimes), result.Length)
	return result, nil
}

//释放结果
func (w *FastHttpWorker) FreeFastHttpResult(res *FastHttpWorkerResponse) {
	if res == nil {
		return
	}
	if res.Response != nil {
		fasthttp.ReleaseResponse(res.Response)
	}
	if res.Request != nil {
		fasthttp.ReleaseRequest(res.Request)
	}

	res.Request = nil
	res.Response = nil
	res.BlockResult = nil
	res = nil
}

//更新结果
func (w *FastHttpWorker) UpdateFastHttpResult(res *FastHttpWorkerResponse) {
	//更新结束时间
	res.BlockResult.EndNanoTime = time.Now().UnixNano()
	res.EndNanoTime = res.BlockResult.EndNanoTime

	//更新块信息
	res.BlockResult.TotalNanoTime = res.BlockResult.UsedNanoTime
	res.BlockResult.UsedMillTimes, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(res.BlockResult.UsedNanoTime)/float64(time.Millisecond)), 64)
	if res.BlockResult.UsedMillTimes < 0 {
		res.BlockResult.UsedMillTimes = 1.0
	}
	res.BlockResult.TotalMillTimes = res.BlockResult.UsedMillTimes
	res.BlockResult.Code = res.Response.StatusCode()
	res.BlockResult.Length = int64(res.Response.Header.ContentLength())
	res.BlockResult.Header = res.Response.Header
	res.BlockResult.ByteSpeed, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(res.BlockResult.Size)/(res.BlockResult.UsedMillTimes/1000)), 64)

	//更新整体信息
	res.TotalNanoTime = time.Now().UnixNano() - res.StartNanoTime
	res.TotalMillTimes, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(res.TotalNanoTime)/float64(time.Millisecond)), 64)
	res.UsedMillTimes, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(res.UsedNanoTime)/float64(time.Millisecond)), 64)
	if res.TotalNanoTime < 0 {
		res.TotalNanoTime = 1.0
	}
	if res.UsedMillTimes < 0 {
		res.UsedMillTimes = 1.0
	}
	res.Code = res.Response.StatusCode()
	res.Length = int64(res.Response.Header.ContentLength())
	res.Header = res.Response.Header
	res.ByteSpeed, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(res.Size)/(res.UsedMillTimes/1000)), 64)
}

//默认FastWorkerResponse
func DefaultFastWorkerResponseFunc(ctx *FastHttpWorkerResponse) {
	return
}