package httpclient

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"time"
)

/*
 * fasthttp请求工具类
 */

// 默认超时时间
const REQ_TIMEOUT = 5 * time.Second

type options struct {
	requestUrl     string            // 请求URL
	requestParams  interface{}       // 请求参数
	requestMethod  string            // 请求方法 POST or GET
	headers        map[string]string // 请求头参数
	requestTimeout time.Duration     // 请求超时时间
}

// 请求参数
type RequestOptions options

// 请求URL
func (opt *RequestOptions) Url(url string) *RequestOptions {
	opt.requestUrl = url
	return opt
}

// 请求参数
func (opt *RequestOptions) Params(params interface{}) *RequestOptions {
	opt.requestParams = params
	return opt
}

// 请求方法
func (opt *RequestOptions) Method(method string) *RequestOptions {
	opt.requestMethod = method
	return opt
}

// 请求头参数
func (opt *RequestOptions) BasicAuth(headers map[string]string) *RequestOptions {
	opt.headers = headers
	return opt
}

// 请求超时时间
func (opt *RequestOptions) Timeout(timeout time.Duration) *RequestOptions {
	opt.requestTimeout = timeout
	return opt
}

/* http请求返回结果 */
func DoHttpExecute(result interface{}, ops *RequestOptions) error {
	return doHttpPostForFast(result, ops)
}

// http请求
func doHttpPostForFast(result interface{}, ops *RequestOptions) error {
	method := ops.requestMethod
	if len(method) <= 0 {
		method = "POST"
	}

	// request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// set header
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(method)
	if ops.headers != nil && len(ops.headers) > 0 {
		for key, val := range ops.headers {
			req.Header.Set(key, val)
		}
	}

	// set url
	req.SetRequestURI(ops.requestUrl)
	// set params
	if ops.requestParams != nil {
		requestBody, err := json.Marshal(ops.requestParams)
		if err != nil {
			return err
		}
		req.SetBody(requestBody)
	}

	// response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// timeout
	timeout := ops.requestTimeout
	if timeout.Nanoseconds() <= 0 {
		timeout = REQ_TIMEOUT
	}

	// do post
	if err := fasthttp.DoTimeout(req, resp, timeout); err != nil {
		return err
	}

	// json format
	err := json.Unmarshal(resp.Body(), result)
	if err != nil {
		return err
	}
	return nil
}
