package httpc

// Author: recallsong
// Email: songruiguo@qq.com

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// 支持的HTTP method
const (
	POST    = "POST"
	GET     = "GET"
	HEAD    = "HEAD"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

// 支持的ContentType类型
const (
	TypeTextHtml          = "text/html"
	TypeTextPlain         = "text/plain"
	TypeApplicationJson   = "application/json"
	TypeApplicationXml    = "application/xml"
	TypeApplicationForm   = "application/x-www-form-urlencoded"
	TypeApplicationStream = "application/octet-stream"
	TypeMultipartFormData = "multipart/form-data"
)

// HttpC 发起http请求的Client
type HttpC struct {
	Context       *Context       // 上下文
	Request       *http.Request  // 请求
	Response      *http.Response // 响应
	BaseURL       string         // 请求url基地址
	URL           string         // 请求的url
	QueryData     url.Values     // 请求url的query参数
	SendMediaType string         // 请求的ContentType
	Data          interface{}    // 要发送的数据体
	Error         error          // 请求发生的错误
	SucStatus     int            // 指定成功的状态码，不匹配则以Error的形式返回
}

// New 创建一个HttpC类型
func New(baseUrl string) *HttpC {
	c := &HttpC{
		Context: DefaultContext,
		Request: &http.Request{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Method:     "GET",
		},
		BaseURL:   baseUrl,
		URL:       baseUrl,
		QueryData: url.Values{},
	}
	return c
}

// Reset 重置HttpC
func (c *HttpC) Reset() *HttpC {
	c.Request = &http.Request{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Method:     "GET",
	}
	c.Response = nil
	c.URL = c.BaseURL
	c.QueryData = url.Values{}
	c.SendMediaType = ""
	c.Data = nil
	c.Error = nil
	c.SucStatus = 0
	return c
}

// SetContext 设置请求上下文
func (c *HttpC) SetContext(context *Context) *HttpC {
	c.Context = context
	return c
}

// Path 追加url路径
func (c *HttpC) Path(path string) *HttpC {
	if c.URL != "" &&
		!strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") || !strings.HasPrefix(path, "unix://") {
		if strings.HasSuffix(c.URL, "/") {
			if strings.HasPrefix(path, "/") {
				c.URL += path[0:]
			} else {
				c.URL += path
			}
		} else {
			if strings.HasPrefix(path, "/") {
				c.URL += path
			} else {
				c.URL += "/" + path
			}
		}
	} else {
		c.URL = path
	}
	return c
}

// EscapedPath 对path进行url编码后，再追加到当前的url上
func (c *HttpC) EscapedPath(path interface{}) *HttpC {
	return c.Path(url.QueryEscape(fmt.Sprint(path)))
}

// Query 增加url的query参数
func (c *HttpC) Query(name, value string) *HttpC {
	c.QueryData.Add(name, value)
	return c
}

// Header 设置http的请求头，会覆盖已存在的请求头
func (c *HttpC) Header(key, value string) *HttpC {
	c.Request.Header.Set(key, value)
	return c
}

// AddHeader 添加http的请求头
func (c *HttpC) AddHeader(key, value string) *HttpC {
	c.Request.Header.Add(key, value)
	return c
}

// ContentType 设置http头中的Content-Type
func (c *HttpC) ContentType(value string) *HttpC {
	c.Request.Header.Set("Content-Type", value)
	return c
}

// BasicAuth 设置基于Basic认证的信息
func (c *HttpC) BasicAuth(userName, password string) *HttpC {
	c.Request.SetBasicAuth(userName, password)
	return c
}

// AddCookie 添加cookie
func (c *HttpC) AddCookie(ck *http.Cookie) *HttpC {
	c.Request.AddCookie(ck)
	return c
}

// Body 设置请求数据体
func (c *HttpC) Body(body interface{}, mediaType ...string) *HttpC {
	if len(mediaType) > 0 {
		c.SendMediaType = mediaType[0]
	} else {
		header := c.Request.Header["Content-Type"]
		if len(header) > 0 {
			c.SendMediaType = header[0]
		} else {
			c.SendMediaType = ""
		}
	}
	c.Data = body
	return c
}

// SuccessStatus 设置成功的状态码
func (c *HttpC) SuccessStatus(sucStatus int) *HttpC {
	c.SucStatus = sucStatus
	return c
}

// Get 发送Get请求，并根据mediaType将结果解析到out中
func (c *HttpC) Get(out interface{}, mediaType ...string) error {
	return c.sendReuqest(GET, out, mediaType...)
}

// Post 发送Post请求，并根据mediaType将结果解析到out中
func (c *HttpC) Post(out interface{}, mediaType ...string) error {
	return c.sendReuqest(POST, out, mediaType...)
}

// Put 发送Put请求，并根据mediaType将结果解析到out中
func (c *HttpC) Put(out interface{}, mediaType ...string) error {
	return c.sendReuqest(PUT, out, mediaType...)
}

// Patch 发送Patch请求，并根据mediaType将结果解析到out中
func (c *HttpC) Patch(out interface{}, mediaType ...string) error {
	return c.sendReuqest(PATCH, out, mediaType...)
}

// Delete 发送Delete请求，并根据mediaType将结果解析到out中
func (c *HttpC) Delete(out interface{}, mediaType ...string) error {
	return c.sendReuqest(DELETE, out, mediaType...)
}

// Options 发送Options请求，并根据mediaType将结果解析到out中
func (c *HttpC) Options(out interface{}, mediaType ...string) error {
	return c.sendReuqest(OPTIONS, out, mediaType...)
}

// throwError 返回error
func (c *HttpC) throwError() error {
	if err := c.Context.invokeCbOnError(c); err != nil {
		return err
	}
	return c.Error
}

// sendReuqest 发送请求
func (c *HttpC) sendReuqest(method string, out interface{}, recvMediaType ...string) error {
	if c.Error != nil {
		return c.throwError()
	}
	c.Request.Method = method
	if c.initReuqest() != nil {
		return c.throwError()
	}
	if c.Request.Body != nil {
		defer c.Request.Body.Close()
	}
	if c.Error = c.Context.invokeCbBeforeSend(c); c.Error != nil {
		return c.throwError()
	}
	interval := c.Context.RetryInterval
	retring := 0
	for retring = 0; c.Context.Retries == -1 || retring <= c.Context.Retries; {
		c.Response, c.Error = c.Context.Client.Do(c.Request)
		if c.Error != nil {
			retring++
			if c.Context.Retries == -1 || retring <= c.Context.Retries {
				if interval > 0 && c.Context.RetryFactor > 0 {
					time.Sleep(interval)
					if c.Context.Retries != -1 {
						interval = time.Duration(float64(interval) * c.Context.RetryFactor)
					}
				}
				if err := c.Context.invokeCbOnRetring(c, retring, interval); err != nil {
					c.Error = err
					return c.throwError()
				}
				continue
			}
		}
		break
	}
	if c.Error != nil {
		return c.throwError()
	}
	defer func() {
		err := recover()
		if err != nil {
			c.Response.Body.Close()
		} else {
			if _, ok := out.(**http.Response); !ok {
				c.Response.Body.Close()
			}
		}
	}()
	if c.Error = c.Context.invokeCbAfterSend(c); c.Error != nil {
		return c.throwError()
	}
	if c.SucStatus != 0 {
		if c.SucStatus < 0 {
			if c.Response.StatusCode > -c.SucStatus {
				c.Error = fmt.Errorf("error http status %d , expect <= %d", c.Response.StatusCode, -c.SucStatus)
				return c.throwError()
			}
		} else {
			if c.Response.StatusCode != c.SucStatus {
				c.Error = fmt.Errorf("error http status %d , expect %d", c.Response.StatusCode, c.SucStatus)
				return c.throwError()
			}
		}
	}
	var mediaType string
	if len(recvMediaType) > 0 {
		mediaType = recvMediaType[0]
	} else {
		mediaType = c.Response.Header.Get("Content-Type")
	}
	c.Error = c.readResponse(c.Response, out, mediaType)
	if c.Error != nil {
		return c.throwError()
	}
	return nil
}

// initReuqest 初始化Request
func (c *HttpC) initReuqest() error {
	//init url
	u, err := url.Parse(c.URL)
	if err != nil {
		c.Error = err
		return c.Error
	}
	//set query params
	q := u.Query()
	for k, v := range c.QueryData {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	u.RawQuery = q.Encode()
	c.Request.URL = u
	// set body
	mediaType := c.SendMediaType
	switch c.Request.Method {
	case POST, PUT, PATCH:
		if c.Data == nil {
			break
		}
		var readerGetter ReqBodyReader
		var find bool
		typ := reflect.TypeOf(c.Data)
		if c.Context.BodyReaders != nil {
			readers := c.Context.BodyReaders
			if readers.ReqBodyTypeReaders != nil {
				readerGetter, find = readers.ReqBodyTypeReaders[typ]
			}
			if find == false {
				if readerGetter, find = GlobalReqBodyTypeReaders[typ]; find == false {
					if mediaType == "" {
						mediaType = strings.TrimSpace(strings.ToLower(strings.Split(c.Request.Header.Get("Content-Type"), ";")[0]))
					} else {
						mediaType = strings.TrimSpace(strings.ToLower(strings.Split(mediaType, ";")[0]))
					}
					if readers.ReqBodyMediaReaders != nil {
						if readerGetter, find = readers.ReqBodyMediaReaders[mediaType]; find == false {
							readerGetter, find = GlobalReqBodyMediaReaders[mediaType]
						}
					} else {
						readerGetter, find = GlobalReqBodyMediaReaders[mediaType]
					}
				}
			}
		} else {
			if readerGetter, find = GlobalReqBodyTypeReaders[typ]; find == false {
				if mediaType == "" {
					mediaType = strings.TrimSpace(strings.ToLower(strings.Split(c.Request.Header.Get("Content-Type"), ";")[0]))
				} else {
					mediaType = strings.TrimSpace(strings.ToLower(strings.Split(mediaType, ";")[0]))
				}
				readerGetter, find = GlobalReqBodyMediaReaders[mediaType]
			}
		}
		if find == false {
			c.Error = ErrReqBodyType
			return c.Error
		}
		c.Request.Body, c.Error = readerGetter(c.Request, typ, c.Data)
	case GET, HEAD, DELETE, OPTIONS:
		break
	default:
		c.Error = fmt.Errorf("invalid http method %s", c.Request.Method)
		return c.Error
	}
	return c.Error
}

// readResponse 从响应中读取数据到out中
func (c *HttpC) readResponse(resp *http.Response, out interface{}, mediaType string) error {
	if out == nil {
		return nil
	}
	var readerGetter RespBodyReader
	var find bool
	typ := reflect.TypeOf(out)
	if c.Context.BodyReaders != nil {
		readers := c.Context.BodyReaders
		if readers.RespBodyTypeReaders != nil {
			readerGetter, find = readers.RespBodyTypeReaders[typ]
		}
		if find == false {
			if readerGetter, find = GlobalRespBodyTypeReaders[typ]; find == false {
				mediaType = strings.TrimSpace(strings.ToLower(strings.Split(mediaType, ";")[0]))
				if readers.RespBodyMediaReaders != nil {
					if readerGetter, find = readers.RespBodyMediaReaders[mediaType]; find == false {
						readerGetter, find = GlobalRespBodyMediaReaders[mediaType]
					}
				} else {
					readerGetter, find = GlobalRespBodyMediaReaders[mediaType]
				}
			}
		}
	} else {
		if readerGetter, find = GlobalRespBodyTypeReaders[typ]; find == false {
			mediaType = strings.TrimSpace(strings.ToLower(strings.Split(mediaType, ";")[0]))
			readerGetter, find = GlobalRespBodyMediaReaders[mediaType]
		}
	}
	if find == false {
		c.Error = ErrRespOutType
		return c.Error
	}
	c.Error = readerGetter(resp, resp.Body, typ, out)
	return c.Error
}

// String 将Httpc转换为字符串
func (c *HttpC) String() string {
	urlStr := c.URL
	if c.Request.URL != nil {
		urlStr = c.Request.URL.String()
	} else {
		u, err := url.Parse(urlStr)
		if err == nil {
			q := u.Query()
			for k, v := range c.QueryData {
				for _, vv := range v {
					q.Add(k, vv)
				}
			}
			u.RawQuery = q.Encode()
			urlStr = u.String()
		}
	}
	if c.Response != nil {
		urlStr += " [ " + c.Response.Status + " ]"
	}
	if c.Error != nil {
		urlStr += " -> " + c.Error.Error()
	}
	return urlStr
}

// DumpRequest 将http请求明细输出为string
func (c *HttpC) DumpRequest(dumpBody ...bool) string {
	if c.Request.URL == nil {
		if c.initReuqest() != nil {
			return fmt.Sprintln("DumpRequest error: ", c.Error.Error())
		}
		if c.Request.Body != nil {
			defer c.Request.Body.Close()
		}
	}
	body := false
	if len(dumpBody) > 0 {
		body = dumpBody[0]
	}
	dump, err := httputil.DumpRequest(c.Request, body)
	if nil != err {
		return fmt.Sprintln("DumpRequest error: ", err.Error())
	}
	return fmt.Sprintf("HTTP Request: \n%s\n", string(dump))
}

// DumpResponse 将http请求明细输出为string
func (c *HttpC) DumpResponse(dumpBody ...bool) string {
	body := false
	if len(dumpBody) > 0 {
		body = dumpBody[0]
	}
	dump, err := httputil.DumpResponse(c.Response, body)
	if nil != err {
		return fmt.Sprintln("DumpResponse error: ", err.Error())
	}
	return fmt.Sprintf("HTTP Response: \n%s\n", string(dump))
}
