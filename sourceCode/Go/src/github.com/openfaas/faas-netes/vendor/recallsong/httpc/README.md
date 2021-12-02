# httpc简介
httpc这是一个发起http请求的客户端库。

它具有的特色包括：简单易用、易于扩展、支持链式调用、支持多种请求和响应格式的处理等。

特别适合用来调用RESTful风格的接口。
    
# 下载
    go get github.com/recallsong/httpc
    
# Api文档
查看 [在线Api文档](https://godoc.org/github.com/recallsong/httpc)

我们也可以利用godoc工具在本地查看api文档：
```
godoc -http=:9090
```
在浏览器中查看地址：

http://localhost:9090/pkg/github.com/recallsong/httpc
    
# 快速入门
## 最简单的使用方式
```go
var resp string
// GET http://localhost/hello?name=RecallSong
err := httpc.New("http://localhost").Path("hello").Query("name", "RecallSong").Get(&resp)
if err != nil {
    fmt.Println(resp) // 以字符串方式获取响应的数据
} else {
    fmt.Println(err)
}
```
## 设置请求头和Cookie等
```go
var resp string
err := httpc.New("http://localhost").Path("/hello").Query("param", "value").
    Header("MyHeader", "HeaderValue").
    AddCookie(&http.Cookie{Name: "cookieName", Value: "cookieValue"}).
    Body("body data").Post(&resp)
if err != nil {
    fmt.Println(resp) // 以字符串方式获取响应的数据
} else {
    fmt.Println(err)
}
```
## 发送和接收json格式的数据
### 使用map传递数据
```go
body := map[string]interface{}{
    "name": "RecallSong",
    "age":  18,
}
var resp map[string]interface{}
// 根据请求的Content-Type自动对数据进行转换
err := httpc.New("http://localhost").Path("json").
    ContentType(httpc.TypeApplicationJson).
    Body(body). // body转变为 {"name":"RecallSong","age":18}
    Post(&resp) // 根据响应中的Content-Type，将返回的数据解析到resp中
fmt.Println(err, resp)

// 如果请求或响应没有指定Content-Type，或是不正确，也可以强制指定转换格式类型
err = httpc.New("http://localhost").Path("json").
    Body(body, httpc.TypeApplicationJson). // body转变为 {"name":"RecallSong","age":18}
    Post(&resp, httpc.TypeApplicationJson) // 将返回的数据按json格式解析到map中
fmt.Println(err, resp)
```
### 使用struct传递数据
```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
body := Person{Name: "RecallSong", Age: 18}
var resp Person
err := httpc.New("http://localhost").Path("json").
    Body(body, httpc.TypeApplicationJson).
    Post(&resp, httpc.TypeApplicationJson)
fmt.Println(err, resp)
```
## 发送和接收xml格式的数据
```go
type Person struct {
    Name string `xml:"name"`
    Age  int    `xml:"age"`
}
body := Person{Name: "RecallSong", Age: 18}
var resp Person
err := httpc.New("http://localhost").Path("xml").
    Body(body, httpc.TypeApplicationXml). // 数据转变为xml格式
    Post(&resp, httpc.TypeApplicationXml)
fmt.Println(err, resp)
``` 
## 发送表单参数
### 使用结构体发送
```go
sbody := struct {
    Name string `form:"name"`
    Age  int    `form:"age"`
}{
    Name: "RecallSong",
    Age:  18,
}
var resp string
err := httpc.New("http://localhost").Path("echo").
    Body(sbody, httpc.TypeApplicationForm). // 将结构体转变为form格式的数据体
    Post(&resp)
fmt.Println(err, resp)
```
### 使用map发送
```go
mbody := map[string]interface{}{
    "name": "RecallSong",
    "age":  19,
}
var resp string
err := httpc.New("http://localhost").Path("echo").
    Body(mbody, httpc.TypeApplicationForm). // 将map变为form格式的数据体
    Post(&resp)
fmt.Println(err, resp)
```
### 使用url.Values发送
```go
ubody := url.Values{}
ubody.Set("name", "RecallSong")
ubody.Set("age", "20")
var resp string
err := httpc.New("http://localhost").Path("echo").
    Body(ubody). // 将url.Values类型转变form格式的数据体
    Post(&resp)
fmt.Println(err, resp)
```
## 自动编码url路径参数
```go
var resp string
// 可以自动编码url路径参数
err := httpc.New("http://localhost").EscapedPath("recall/Song").EscapedPath(18).Get(&resp)
// 请求地址为 http://localhost/recall%2FSong/18
fmt.Println(err, resp)
```
## 上传文件
### 方式1
```go
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
err = httpc.New("http://localhost").Path("echo").
    Body(body, httpc.TypeMultipartFormData).Post(&resp)
fmt.Println(err)
```
### 方式2
```go
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
err = httpc.New("http://localhost").Path("echo").
    Body(body, httpc.TypeMultipartFormData).Post(&resp)
fmt.Println(err)
```
## 接收响应数据
```go
// 前面的例子我们知道了可以接收json和xml格式的数据，也可以接收数据到一个string变量中
// 除此之外，我们还可以有一下几种方式接收数据

// []byte 方式接收
var bytesResp []byte
err := httpc.New("http://localhost").Path("hello").Get(&bytesResp)
fmt.Println(err, bytesResp)

// *http.Response 方式接收
var resp *http.Response
err := httpc.New("http://localhost").Path("hello").Get(&resp)
if err != nil {
    fmt.Println(err)
} else {
    // 注意这种方式要关闭Body
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    fmt.Println(err, string(body))
}
```
## 下载文件
### 方式1
```go
// 默认方式保存文件
err := httpc.New("http://localhost").Path("echo").Body("content").Post(httpc.FilePath("download1.txt"))
fmt.Println(err)
```
### 方式2
```go
err := httpc.New("http://localhost").Path("echo").Body("content").Post(&httpc.SaveInfo{
    Path:     "download2.txt",
    Override: true,
    Mode:     0777})
fmt.Println(err)
```
## 指定成功的http状态码
```go
// 如果返回的状态码与指定的状态码不匹配，则返回一个error
err := httpc.New("http://localhost").Path("not_exist").
    SuccessStatus(200).Get(nil)
fmt.Println(err)
// Output:
// error http status 404 , expect 200
```
## 请求上下文
```go
// 请求上下文中包含了每次请求的设置、连接设置等，所有请求应该尽量共享Context
// 我们可以设置回调通知的函数
ctx := httpc.NewContext().
    AddCbBeforeSend(func(client *httpc.HttpC, args ...interface{}) error {
        fmt.Println("before request")
        return nil
    }).
    AddCbAfterSend(func(client *httpc.HttpC, args ...interface{}) error {
        fmt.Println("after response")
        return nil
    }).
    AddCbOnError(func(client *httpc.HttpC, args ...interface{}) error {
        fmt.Println("on error")
        return nil
    }).
    SetConnectReadTimeout(30*time.Second, 30*time.Second)
var resp string
err := httpc.New("http://localhost").Path("hello").SetContext(ctx).Get(&resp)
fmt.Println(err, resp)

// 库默认生成了一个上下文实例 httpc.DefaultContext，它并没有加锁保护，所以尽量在所有请求前设置好它
// 改变httpc.DefaultContext会影响所有未调用过SetContext的请求
httpc.DefaultContext.SetConnectReadTimeout(30*time.Second, 30*time.Second)
err = httpc.New("http://localhost").Path("hello").Get(&resp)
fmt.Println(err, resp)
```
## 超时设置
```go
err := httpc.New("http://localhost").Path("timeout").
    SetContext(httpc.NewContext().SetConnectReadTimeout(time.Second, time.Second)).
    Get(nil)
fmt.Println(err)
```
## 请求重试
```go
err := httpc.New("http://not_exist/").Path("not_exist").
    SetContext(httpc.NewContext().AddCbOnRetring(func(c *httpc.HttpC, args ...interface{}) error {
        fmt.Printf("retring %v, next interval %v\n", args[0], args[1])
        return nil
    }).SetRetryConfig(3, time.Second, 2)). // 重试3次，重试时间间隔依次为：2s, 4s, 8s
    Get(nil)
fmt.Println(err)

// Output:
// retring 1, next interval 2s
// retring 2, next interval 4s
// retring 3, next interval 8s
// Get http://not_exist/not_exist: dial tcp: lookup not_exist: no such host
```
## 自定义请求或响应处理器
```go
// httpc库已经注册了一些通用的请求和响应处理器，但我们也可以额外添加处理器
ctx := httpc.NewContext()
ctx.BodyReaders = httpc.NewBodyReaders()
ctx.BodyReaders.RespBodyTypeReaders[reflect.TypeOf((*int)(nil))] = func(resp *http.Response, reader io.ReadCloser, typ reflect.Type, out interface{}) error {
    output := out.(*int)
    *output = resp.StatusCode
    return nil
}
// 返回响应状态码
var status int
err := httpc.New("http://localhost").Path("hello").
    SetContext(ctx).
    Get(&status)
fmt.Println(err, status)
// Output:
// <nil> 200
```
## 其他特性
    请参考Api文档

# License
[MIT](https://github.com/recallsong/httpc/blob/master/LICENSE)
