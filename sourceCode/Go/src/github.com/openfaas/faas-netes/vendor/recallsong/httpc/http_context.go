package httpc

// Author: recallsong
// Email: songruiguo@qq.com

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"time"
)

// DefaultContext 默认的请求上下文
var DefaultContext = NewContext()

// Callback 以回调方式通知内部所发生的事件
type Callback func(c *HttpC, args ...interface{}) error

// BodyReaders 请求体读取器
type BodyReaders struct {
	ReqBodyTypeReaders   map[reflect.Type]ReqBodyReader  // 根据类型获取reader
	ReqBodyMediaReaders  map[string]ReqBodyReader        // 根据MediaType获取reader
	RespBodyTypeReaders  map[reflect.Type]RespBodyReader // 根据类型获取reader
	RespBodyMediaReaders map[string]RespBodyReader       // 根据MediaType获取reader
}

// NewBodyReaders 创建一个BodyReaders
func NewBodyReaders() *BodyReaders {
	return &BodyReaders{
		ReqBodyTypeReaders:   make(map[reflect.Type]ReqBodyReader),
		ReqBodyMediaReaders:  make(map[string]ReqBodyReader),
		RespBodyTypeReaders:  make(map[reflect.Type]RespBodyReader),
		RespBodyMediaReaders: make(map[string]RespBodyReader),
	}
}

// Context http请求上下文，管理所有请求公用的对象
type Context struct {
	Client      *http.Client // 请求的客户端
	BodyReaders *BodyReaders // 获取请求和响应中的body的reader

	CbBeforeSend  []Callback    // 在发送请求前调用
	CbAfterSend   []Callback    // 在发送请求后调用
	CbOnError     []Callback    // 在发生错误时调用
	CbOnRetring   []Callback    // 在请求重试时调用
	Retries       int           // 重试次数，-1表示一直重试
	RetryInterval time.Duration // 重试间隔
	RetryFactor   float64       // 重试因子，影响每次重试的时间间隔
}

// NewContext 创建一个Context实例
func NewContext() *Context {
	return &Context{
		Client:      &http.Client{Transport: &http.Transport{}},
		RetryFactor: 1,
	}
}

// AddCbBeforeSend 添加发送请求前的通知回调函数
func (c *Context) AddCbBeforeSend(cb Callback) *Context {
	if c.CbBeforeSend == nil {
		c.CbBeforeSend = make([]Callback, 1)
		c.CbBeforeSend[0] = cb
	} else {
		c.CbBeforeSend = append(c.CbBeforeSend, cb)
	}
	return c
}

// AddCbAfterSend 添加发送请求后的通知回调函数
func (c *Context) AddCbAfterSend(cb Callback) *Context {
	if c.CbAfterSend == nil {
		c.CbAfterSend = make([]Callback, 1)
		c.CbAfterSend[0] = cb
	} else {
		c.CbAfterSend = append(c.CbAfterSend, cb)
	}
	return c
}

// AddCbOnError 添加发生错误时的通知回调函数
func (c *Context) AddCbOnError(cb Callback) *Context {
	if c.CbOnError == nil {
		c.CbOnError = make([]Callback, 1)
		c.CbOnError[0] = cb
	} else {
		c.CbOnError = append(c.CbOnError, cb)
	}
	return c
}

// AddCbOnRetring 添加请求重试的通知回调函数
func (c *Context) AddCbOnRetring(cb Callback) *Context {
	if c.CbOnRetring == nil {
		c.CbOnRetring = make([]Callback, 1)
		c.CbOnRetring[0] = cb
	} else {
		c.CbOnRetring = append(c.CbOnRetring, cb)
	}
	return c
}

// invokeCbBeforeSend 调用所有CbBeforeSend回调函数
func (c *Context) invokeCbBeforeSend(hc *HttpC, args ...interface{}) error {
	if c.CbBeforeSend != nil {
		length := len(c.CbBeforeSend)
		for i := length - 1; i >= 0; i-- {
			err := c.CbBeforeSend[i](hc, args...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// invokeCbAfterSend 调用所有CbAfterSend回调函数
func (c *Context) invokeCbAfterSend(hc *HttpC, args ...interface{}) error {
	if c.CbAfterSend != nil {
		length := len(c.CbAfterSend)
		for i := length - 1; i >= 0; i-- {
			err := c.CbAfterSend[i](hc, args...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// invokeCbOnError 调用所有CbOnError回调函数
func (c *Context) invokeCbOnError(hc *HttpC, args ...interface{}) error {
	if c.CbOnError != nil {
		length := len(c.CbOnError)
		for i := length - 1; i >= 0; i-- {
			err := c.CbOnError[i](hc, args...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// invokeCbOnRetring 调用所有CbOnError回调函数
func (c *Context) invokeCbOnRetring(hc *HttpC, args ...interface{}) error {
	if c.CbOnRetring != nil {
		length := len(c.CbOnRetring)
		for i := length - 1; i >= 0; i-- {
			err := c.CbOnRetring[i](hc, args...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// getTransport 从Client中获取*http.Transport
func (c *Context) getTransport() *http.Transport {
	if c.Client.Transport == nil {
		c.Client.Transport = &http.Transport{}
	}
	if t, ok := c.Client.Transport.(*http.Transport); ok {
		return t
	}
	c.Client.Transport = &http.Transport{}
	return c.Client.Transport.(*http.Transport)
}

// SetClient 设置http.Client
func (c *Context) SetClient(client *http.Client) *Context {
	c.Client = client
	return c
}

// SetTotalTimeout 设置总超时，包括连接、所有从定向、读数据的时间，直到数据被读完
func (c *Context) SetTotalTimeout(timeout time.Duration) *Context {
	c.Client.Timeout = timeout
	return c
}

// SetConnectReadTimeout 设置请求的连接超时和读数据超时
func (c *Context) SetConnectReadTimeout(connectTimeout time.Duration, readTimeout time.Duration) *Context {
	c.getTransport().Dial = func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(readTimeout))
		return conn, nil
	}
	return c
}

// SetRetryConfig 设置重试的配置
func (c *Context) SetRetryConfig(retries int, interval time.Duration, factor float64) *Context {
	c.Retries = retries
	c.RetryInterval = interval
	c.RetryFactor = factor
	return c
}

// EnableCookie 启用CookieJar
func (c *Context) EnableCookie(enable bool) *Context {
	if enable {
		c.Client.Jar = nil
	} else {
		if c.Client.Jar == nil {
			jar, _ := cookiejar.New(nil)
			c.Client.Jar = jar
		}
	}
	return c
}

// SetTLSClientConfig 设置TLSClientConfig
func (c *Context) SetTLSClientConfig(config *tls.Config) *Context {
	c.getTransport().TLSClientConfig = config
	return c
}

// SetCheckRedirect 设置CheckRedirect
func (c *Context) SetCheckRedirect(cr func(req *http.Request, via []*http.Request) error) *Context {
	c.Client.CheckRedirect = cr
	return c
}

// SetProxy 设置请求代理
func (c *Context) SetProxy(proxy func(*http.Request) (*url.URL, error)) *Context {
	c.getTransport().Proxy = proxy
	return c
}

// Copy 复制一份
func (c *Context) Copy() *Context {
	return &Context{
		Client:        c.Client,
		BodyReaders:   c.BodyReaders,
		CbBeforeSend:  c.CbBeforeSend,
		CbAfterSend:   c.CbAfterSend,
		CbOnError:     c.CbOnError,
		CbOnRetring:   c.CbOnRetring,
		Retries:       c.Retries,
		RetryInterval: c.RetryInterval,
		RetryFactor:   c.RetryFactor,
	}
}
