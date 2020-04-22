package boomer_client

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"falcon/slave/site/boomer"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)


//请求配置
type HttpRequestConfig struct {
	Url     		string
	Route   		string
	ContentType		string
	Method 			string
	Params  		map[string]string
	Headers 		map[string]string
	Body    		[]byte
	Request 		*http.Request
	Blocks			int
	BlocksWait		int64
	NoRespBody		bool
	Close 			bool
	GetMd5			bool
	WorkerResponse 	func(ctx *HttpWorkerResponse)
}

//返回值结构
type HttpWorkerResponse struct {
	Request 		*http.Request
	Response 		*http.Response
	BlockResult		*HttpResponseResult
	*HttpResponseResult
}

type HttpResponseResult struct {
	Body   			[]byte
	Code   			int
	Length 			int64
	Size 			int64
	Header 			http.Header
	Proto  			string
	TLS    			*tls.ConnectionState
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

type HttpWorker struct {
	Client *HttpWorkerClient
}

type HttpWorkerClient struct {
	*http.Client
}

var DefaultHttpBlocks = 1024*1024


//创建Client
func NewHttpClient(c ...*HttpWorkerClient) *HttpWorker {
	defaultClient :=  &HttpWorkerClient{
		Client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (
					net.Conn, error) {
					c, err := net.DialTimeout(network, addr, time.Second*60)
					if err != nil {
						return nil, err
					}
					return c, nil
				},
				ResponseHeaderTimeout: time.Second * 60,
				MaxIdleConnsPerHost:   50000,
				MaxIdleConns:          50000,
				IdleConnTimeout:       time.Second * 60,
			},
		},
	}

	var client HttpWorkerClient
	if len(c) == 0 {
		client = *defaultClient
	} else {
		client = *c[0]
	}
	return &HttpWorker{Client: &client}
}

//创建默认带Http代理的Client
func NewDefaultProxyHttpClient(proxyAddr string) *HttpWorker {
	defaultClient :=  &HttpWorkerClient{
		Client: &http.Client{
			Transport: &http.Transport{
				Proxy: func(_ *http.Request) (*url.URL, error) {
					return url.Parse("http://" + proxyAddr)
				},
				DialContext: func(ctx context.Context, network, addr string) (
					net.Conn, error) {
					c, err := net.DialTimeout(network, addr, time.Second*60)
					if err != nil {
						return nil, err
					}
					return c, nil
				},
				ResponseHeaderTimeout: time.Second * 60,
				MaxIdleConnsPerHost:   50000,
				MaxIdleConns:          50000,
				IdleConnTimeout:       time.Second * 60,
			},
		},
	}
	return &HttpWorker{Client: defaultClient}
}

//发送Post请求, isCatch=False,自动向Master发送事件
func (w *HttpWorker) Post(resConf *HttpRequestConfig, isCatch ...bool) (*HttpWorkerResponse, error) {
	return w.Do("POST", resConf, isCatch...)
}

//发送Get请求
func (w *HttpWorker) Get(resConf *HttpRequestConfig, isCatch ...bool) (*HttpWorkerResponse, error) {
	return w.Do("GET", resConf, isCatch...)
}

//发送Put请求
func (w *HttpWorker) Put(resConf *HttpRequestConfig, isCatch ...bool) (*HttpWorkerResponse, error) {
	return w.Do("PUT", resConf, isCatch...)
}

//发送Delete请求
func (w *HttpWorker) Delete(resConf *HttpRequestConfig, isCatch ...bool) (*HttpWorkerResponse, error) {
	return w.Do("DELETE", resConf, isCatch...)
}

//发送请求
func (w *HttpWorker) sendHttpRequest(method string, resConf *HttpRequestConfig) (*HttpWorkerResponse, error) {
	result := &HttpWorkerResponse{
		Request: nil,
		Response: nil,
		BlockResult: &HttpResponseResult{},
		HttpResponseResult: &HttpResponseResult{},
	}

	urlPath := strings.Trim(resConf.Url, "/") + resConf.Route
	split := "?"
	for k, v := range resConf.Params {
		urlPath += fmt.Sprintf("%s%s=%s", split, k, v)
		split = "&"
	}
	//var req *http.Request
	if resConf.Request == nil {
		var err error
		result.Request, err = http.NewRequest(method, urlPath, bytes.NewReader(resConf.Body))
		if err != nil {
			return nil, err
		}

		for k, v := range resConf.Headers {
			result.Request.Header.Set(k, v)
		}

		if resConf.ContentType != "" {
			result.Request.Header.Set("Content-Type", resConf.ContentType)
		}
	} else {
		result.Request = resConf.Request
	}
	if resConf.Blocks <= 0 {
		resConf.Blocks = DefaultHttpBlocks
	}
	if resConf.WorkerResponse == nil {
		resConf.WorkerResponse = DefaultWorkerResponseFunc
	}
	result.Request.Close = resConf.Close
	result.StartNanoTime = time.Now().UnixNano()
	result.BlockResult.StartNanoTime = result.StartNanoTime

	var doErr error
	result.Response, doErr = w.Client.Do(result.Request)
	result.EndNanoTime = time.Now().UnixNano()
	result.BlockResult.EndNanoTime = result.EndNanoTime
	result.UsedNanoTime += (time.Now().UnixNano() - result.StartNanoTime)
	if doErr != nil {
		return nil, doErr
	}
	defer result.Response.Body.Close()

	var blockErr error
	var allBuf bytes.Buffer
	var bufSize int
	var bufNum int
	md5AllBuf := md5.New()
	for {
		md5Buf := md5.New()
		buf := make([]byte, resConf.Blocks)
		bufNum = 0
		bufSize = 0
		result.BlockResult.StartNanoTime = time.Now().UnixNano()
		for {
			bufNum, blockErr = result.Response.Body.Read(buf[bufSize:])
			bufSize += bufNum
			if blockErr != nil || bufNum == 0 || bufSize == resConf.Blocks {
				break
			}
		}
		result.BlockResult.UsedNanoTime = time.Now().UnixNano() - result.BlockResult.StartNanoTime
		result.UsedNanoTime += result.BlockResult.UsedNanoTime
		if bufSize > 0 {
			if resConf.GetMd5 {
				md5Buf.Write(buf[0:bufSize])
				md5AllBuf.Write(buf[0:bufSize])
				result.BlockResult.Md5 = strings.ToUpper(fmt.Sprintf("%x",md5Buf.Sum(nil)))
				result.Md5 = strings.ToUpper(fmt.Sprintf("%x",md5AllBuf.Sum(nil)))
			}
			if !resConf.NoRespBody {
				result.BlockResult.Body = buf[0:bufSize]
				allBuf.Write(buf[0:bufSize])
				result.Body = allBuf.Bytes()
			}
			buf = nil
		}
		result.Size += int64(bufSize)
		result.BlockResult.Size = int64(bufSize)
		time.Sleep(time.Duration(resConf.BlocksWait)*time.Millisecond)
		w.UpdateHttpResult(result)

		if blockErr != nil {
			if blockErr.Error() != "EOF" {
				result.Error = blockErr
				result.BlockResult.Error = blockErr
			} else {
				blockErr = nil
			}
			resConf.WorkerResponse(result)
			break
		}
		resConf.WorkerResponse(result)
	}
	w.UpdateHttpResult(result)
	return result, blockErr
}

//创建请求(POST/GET/PUT/DELETE)
func (w *HttpWorker) Do(method string, res *HttpRequestConfig, isCatch ...bool) (*HttpWorkerResponse, error) {
	boom := boomer.Slave.Boom
	param := boomer.Param
	minWait := param.MinWait
	maxWait := param.MaxWait
	defer boomer.RandomWait(time.Duration(minWait)*time.Millisecond, time.Duration(maxWait)*time.Millisecond)

	var result *HttpWorkerResponse
	var err error

	result, err = w.sendHttpRequest(method, res)

	if len(isCatch) > 0 {
		if isCatch[0] {
			return result, err
		}
	}

	defer w.FreeHttpResult(result)
	if err != nil {
		boom.RecordFailure(method, res.Route, 0, err.Error())
		return result, err
	}

	if (result.Code < 200) || (result.Code >= 300) {
		err := errors.New(fmt.Sprintf("Response code is %d", result.Code))
		boom.RecordFailure(method, res.Route, int64(result.UsedMillTimes), err.Error())
		return result, err
	}
	boom.RecordSuccess(method, res.Route, int64(result.UsedMillTimes), result.Size)
	return result, nil
}

//释放结果
func (w *HttpWorker) FreeHttpResult(res *HttpWorkerResponse) {
	if res == nil {
		return
	}
	if res.Response != nil {
		_ = res.Response.Body.Close()
	}
	res.Request = nil
	res.Response = nil
	res.BlockResult = nil
	res = nil
}

//更新结果
func (w *HttpWorker) UpdateHttpResult(res *HttpWorkerResponse) {
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
	res.BlockResult.Code = res.Response.StatusCode
	res.BlockResult.Length = res.Response.ContentLength
	res.BlockResult.Header = res.Response.Header
	res.BlockResult.Proto = res.Response.Proto
	res.BlockResult.TLS = res.Response.TLS
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
	res.Code = res.Response.StatusCode
	res.Length = res.Response.ContentLength
	res.Header = res.Response.Header
	res.Proto = res.Response.Proto
	res.TLS = res.Response.TLS
	res.ByteSpeed, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(res.Size)/(res.UsedMillTimes/1000)), 64)
}

//默认WorkerResponse
func DefaultWorkerResponseFunc(ctx *HttpWorkerResponse) {
	return
}