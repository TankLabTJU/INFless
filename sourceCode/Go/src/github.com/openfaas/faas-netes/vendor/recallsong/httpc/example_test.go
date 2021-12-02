package httpc_test

// Author: recallsong
// Email: songruiguo@qq.com

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/recallsong/httpc"
)

type Handler struct {
	Method string
	Func   func(w http.ResponseWriter, r *http.Request)
}

var server *httptest.Server

func startServer() string {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := routers[r.URL.Path]; ok {
			if handler.Method != "" && handler.Method != r.Method {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			handler.Func(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}))
	return server.URL
}

func stopServer() {
	server.Close()
}

var routers = map[string]Handler{
	"/hello": Handler{
		Method: "GET",
		Func: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, "+r.URL.Query().Get("name"))
		},
	},
	"/dump": Handler{
		Func: func(w http.ResponseWriter, r *http.Request) {
			bytes, err := httputil.DumpRequest(r, true)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(w, string(bytes))
		},
	},
	"/echo": Handler{
		Method: "POST",
		Func: func(w http.ResponseWriter, r *http.Request) {
			bytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Write(bytes)
		},
	},
	"/json": Handler{
		Method: "POST",
		Func: func(w http.ResponseWriter, r *http.Request) {
			bytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json; charset=UTF-8")
			w.Write(bytes)
		},
	},
	"/xml": Handler{
		Method: "POST",
		Func: func(w http.ResponseWriter, r *http.Request) {
			bytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/xml; charset=UTF-8")
			w.Write(bytes)
		},
	},
	"/error": Handler{
		Func: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		},
	},
	"/timeout": Handler{
		Func: func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			fmt.Fprintln(w, "Hello")
		},
	},
	"/gzip": Handler{
		Func: func(w http.ResponseWriter, r *http.Request) {
			var buffer bytes.Buffer
			gw := gzip.NewWriter(&buffer)
			fmt.Fprintln(gw, "Hello Hello Hello Hello")
			gw.Close()
			byt := buffer.Bytes()
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.Header().Add("Content-Length", fmt.Sprint(len(byt)))
			w.Header().Add("Content-Encoding", "gzip")
			w.Write(byt)
		},
	},
}

func ExampleHttpC_hello() {
	baseUrl := startServer()
	defer stopServer()
	var resp string
	// 请求 {baseUrl}/hello?name=RecallSong, 并将返回数据读入到resp变量
	err := httpc.New(baseUrl).Path("hello").Query("name", "RecallSong").Get(&resp)
	fmt.Println(err)
	fmt.Println(resp)
	// Output:
	// <nil>
	// Hello, RecallSong
}

func ExampleHttpC_response() {
	baseUrl := startServer()
	defer stopServer()
	var resp *http.Response
	// 请求 {baseUrl}/hello?name=RecallSong, 并将返回数据读入到resp变量
	err := httpc.New(baseUrl).Path("hello").Query("name", "RecallSong").Get(&resp)
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(resp.Status)
		fmt.Println(err, string(body))
	}

	// Output:
	// 200 OK
	// <nil> Hello, RecallSong
}

func ExampleHttpC_dump() {
	baseUrl := startServer()
	defer stopServer()
	var resp []byte
	// 构建更复杂的请求
	err := httpc.New(baseUrl).Path("/dump").Query("param", "value").
		Header("MyHeader", "HeaderValue").
		AddCookie(&http.Cookie{Name: "cook", Value: "testcook"}).
		Body("body data").Post(&resp)
	fmt.Println(err)
	// fmt.Println(string(resp))

	// Output:
	// <nil>
}

func ExampleHttpC_mapBody() {
	baseUrl := startServer()
	defer stopServer()
	body := map[string]interface{}{
		"name": "RecallSong",
		"age":  18,
	}
	var resp map[string]interface{}
	// 根据请求的Content-Type自动对数据进行转换
	err := httpc.New(baseUrl).Path("json").
		ContentType(httpc.TypeApplicationJson).
		Body(body). // body转变为 {"name":"RecallSong","age":18}
		Post(&resp) // 根据响应中的Content-Type，将返回的数据解析到resp中
	fmt.Println(err)
	fmt.Println(resp["name"], resp["age"])

	// 如果请求或响应没有指定Content-Type，或是错误的Content-Type，也可以强制指定转换格式类型
	err = httpc.New(baseUrl).Path("json").
		Body(body, httpc.TypeApplicationJson). // body转变为 {"name":"RecallSong","age":18}
		Post(&resp, httpc.TypeApplicationJson) // 将返回的数据按json格式解析到map中
	fmt.Println(err)
	fmt.Println(resp["name"], resp["age"])
	// Output:
	// <nil>
	// RecallSong 18
	// <nil>
	// RecallSong 18
}

type Person struct {
	Name string `json:"name" form:"name" xml:"name"`
	Age  int    `json:"age" form:"age" xml:"age"`
}

func ExampleHttpC_structBody() {
	baseUrl := startServer()
	defer stopServer()
	body := Person{Name: "RecallSong", Age: 18}
	var resp Person
	// 可以使用结构体来传递数据
	err := httpc.New(baseUrl).Path("echo").
		Body(body, httpc.TypeApplicationJson).
		Post(&resp, httpc.TypeApplicationJson)
	fmt.Println(err)
	fmt.Println(resp.Name, resp.Age)
	// Output:
	// <nil>
	// RecallSong 18
}

func ExampleHttpC_xmlBody() {
	baseUrl := startServer()
	defer stopServer()
	body := Person{Name: "RecallSong", Age: 18}
	var resp Person
	// 发送和接收xml数据
	err := httpc.New(baseUrl).Path("xml").
		Body(body, httpc.TypeApplicationXml). // 数据转变为xml格式
		Post(&resp)
	fmt.Println(err)
	fmt.Println(resp)
	// Output:
	// <nil>
	// {RecallSong 18}
}

func ExampleHttpC_formBody() {
	baseUrl := startServer()
	defer stopServer()
	// struct body
	sbody := struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}{
		Name: "RecallSong",
		Age:  18,
	}
	var resp string
	// 发送form参数
	err := httpc.New(baseUrl).Path("echo").
		Body(sbody, httpc.TypeApplicationForm). // 将结构体转变为form格式的数据体
		Post(&resp)
	fmt.Println(err)
	fmt.Println(resp)

	// map body
	mbody := map[string]interface{}{
		"name": "RecallSong",
		"age":  19,
	}
	err = httpc.New(baseUrl).Path("echo").
		Body(mbody, httpc.TypeApplicationForm). // 将map变为form格式的数据体
		Post(&resp)
	fmt.Println(err)
	fmt.Println(resp)

	// url.Values body
	ubody := url.Values{}
	ubody.Set("name", "RecallSong")
	ubody.Set("age", "20")
	err = httpc.New(baseUrl).Path("echo").
		Body(ubody). // 将url.Values类型转变form格式的数据体
		Post(&resp)
	fmt.Println(err)
	fmt.Println(resp)
	// Output:
	// <nil>
	// age=18&name=RecallSong
	// <nil>
	// age=19&name=RecallSong
	// <nil>
	// age=20&name=RecallSong
}

func ExampleHttpC_error() {
	baseUrl := startServer()
	defer stopServer()
	err := httpc.New(baseUrl).Path("not_exist").
		SetContext(httpc.NewContext().AddCbOnError(func(client *httpc.HttpC, args ...interface{}) error {
			fmt.Println("on error: ", client.Error)
			return nil
		})).
		SuccessStatus(200).Get(nil)
	fmt.Println(err)
	// Output:
	// on error:  error http status 404 , expect 200
	// error http status 404 , expect 200
}

func ExampleHttpC_path() {
	baseUrl := startServer()
	defer stopServer()
	req := httpc.New(baseUrl).EscapedPath("recall/song").EscapedPath(18).DumpRequest()
	fmt.Println(strings.Contains(req, "/recall%2Fsong/18"))
	// Output:
	// true
}

func ExampleHttpC_context() {
	baseUrl := startServer()
	defer stopServer()
	ctx := httpc.NewContext().
		AddCbBeforeSend(func(client *httpc.HttpC, args ...interface{}) error {
			// fmt.Println(client.DumpRequest())
			fmt.Println("before request")
			return nil
		}).
		AddCbAfterSend(func(client *httpc.HttpC, args ...interface{}) error {
			// fmt.Println(client.DumpResponse())
			fmt.Println("after response")
			return nil
		}).
		AddCbOnError(func(client *httpc.HttpC, args ...interface{}) error {
			// fmt.Println(client.Error)
			fmt.Println("on error")
			return nil
		}).
		SetConnectReadTimeout(30*time.Second, 30*time.Second)
	var resp string
	err := httpc.New(baseUrl).Path("hello").Query("name", "Song").SetContext(ctx).SuccessStatus(200).Get(&resp)
	fmt.Println(err, resp)
	// Output:
	// before request
	// after response
	// <nil> Hello, Song
}

func ExampleHttpC_mutipart_map() {
	baseUrl := startServer()
	defer stopServer()
	file, err := os.Open("doc.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	body := map[string]interface{}{
		"file":  file,
		"name":  "RecallSong",
		"age":   18,
		"file2": httpc.FilePath("doc.go:hello.go"), //上传doc.go文件，参数名为file2，文件名为hello.go
	}
	var resp string
	err = httpc.New(baseUrl).Path("echo").
		Body(body, httpc.TypeMultipartFormData).Post(&resp)
	fmt.Println(err)
	// fmt.Println(resp)

	// Output:
	// <nil>
}

func ExampleHttpC_mutipart_struct() {
	baseUrl := startServer()
	defer stopServer()
	file, err := os.Open("doc.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	body := struct {
		Name    string         `form:"name"`
		Address []string       `form:"address"`
		Age     int            `form:"age"`
		File    *os.File       `form:"file" file:"hello.go"`
		File2   httpc.FilePath `form:"file2"`
	}{
		Name:    "RecallSong",
		Address: []string{"HangZhou", "WenZhou"},
		Age:     18,
		File:    file,
		File2:   httpc.FilePath("doc.go:hello2.go"), //上传doc.go文件，参数名为file2，文件名为hello2.go
	}
	var resp string
	err = httpc.New(baseUrl).Path("echo").
		Body(body, httpc.TypeMultipartFormData).Post(&resp)

	fmt.Println(err)
	// fmt.Println(resp)

	// Output:
	// <nil>
}

func ExampleHttpC_download() {
	baseUrl := startServer()
	defer stopServer()
	// 默认方式保存文件
	err := httpc.New(baseUrl).Body("xxx").Path("echo").Post(httpc.FilePath("download1.txt"))
	fmt.Println(err)
	_, err = os.Stat("download1.txt")
	if os.IsNotExist(err) {
		fmt.Println(err)
	}
	// 保存文件的另一种方式
	err = httpc.New(baseUrl).Body("zzz").Path("echo").Post(&httpc.SaveInfo{
		Path:     "download2.txt",
		Override: true,
		Mode:     0777})
	fmt.Println(err)
	_, err = os.Stat("download2.txt")
	if os.IsNotExist(err) {
		fmt.Println(err)
	}
}

func ExampleHttpC_timeout() {
	baseUrl := startServer()
	defer stopServer()
	// 测试读数据超时
	err := httpc.New(baseUrl).Path("timeout").
		SetContext(httpc.NewContext().SetConnectReadTimeout(time.Second, 1*time.Second)).
		Get(nil)
	fmt.Println(err != nil)
	// Output:
	// true
}

func ExampleHttpC_retry() {
	// 测试重试请求
	err := httpc.New("http://not_exist/").Path("not_exist").
		SetContext(httpc.NewContext().AddCbOnRetring(func(c *httpc.HttpC, args ...interface{}) error {
			fmt.Printf("retring %v, next interval %v\n", args[0], args[1])
			return nil
		}).SetRetryConfig(3, time.Second, 2)).
		Get(nil)
	fmt.Println(err)

	// Output:
	// retring 1, next interval 2s
	// retring 2, next interval 4s
	// retring 3, next interval 8s
	// Get http://not_exist/not_exist: dial tcp: lookup not_exist: no such host
}

func ExampleHttpC_gzip() {
	baseUrl := startServer()
	defer stopServer()
	// 测试重试请求
	var resp string
	err := httpc.New(baseUrl).Path("gzip").Get(&resp)
	fmt.Println(err)
	fmt.Println(resp)

	// Output:
	// <nil>
	// Hello Hello Hello Hello
}

func ExampleHttpC_body_reader() {
	baseUrl := startServer()
	defer stopServer()
	ctx := httpc.NewContext()
	ctx.BodyReaders = httpc.NewBodyReaders()
	ctx.BodyReaders.RespBodyTypeReaders[reflect.TypeOf((*int)(nil))] = func(resp *http.Response, reader io.ReadCloser, typ reflect.Type, out interface{}) error {
		output := out.(*int)
		*output = resp.StatusCode
		return nil
	}
	// 返回响应状态码
	var status int
	err := httpc.New(baseUrl).Path("hello").
		SetContext(ctx).
		Get(&status)
	fmt.Println(err, status)

	// Output:
	// <nil> 200
}
