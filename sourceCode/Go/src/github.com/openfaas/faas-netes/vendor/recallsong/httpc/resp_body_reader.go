package httpc

// Author: recallsong
// Email: songruiguo@qq.com

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

// RespBodyReader 根据类型将数据读取到out中
type RespBodyReader func(resp *http.Response, reader io.ReadCloser, typ reflect.Type, out interface{}) error

var (
	// ErrRespOutType 接收响应Body的数据类型错误
	ErrRespOutType = errors.New("invalid response output type")
)

// FilePath 保存文件的路径
type FilePath string

// SaveInfo 保存文件的设置
type SaveInfo struct {
	Path     string
	Mode     os.FileMode
	Override bool
}

// GlobalRespBodyTypeReaders 默认的根据类型获取ResponseBody的Reader映射
var GlobalRespBodyTypeReaders = map[reflect.Type]RespBodyReader{
	reflect.TypeOf((*string)(nil)):         readAllRespWrapper(stringRespBodyReader),
	reflect.TypeOf((*[]byte)(nil)):         readAllRespWrapper(bytesRespBodyReader),
	reflect.TypeOf((**http.Response)(nil)): responseReader,
	reflect.TypeOf((FilePath)("")):         downloadRespBodyReader,
	reflect.TypeOf((*SaveInfo)(nil)):       downloadRespBodyReader,
	reflect.TypeOf(SaveInfo{}):             downloadRespBodyReader,
}

// GlobalRespBodyMediaReaders 默认的根据MediaType获取ResponseBody的Reader映射
var GlobalRespBodyMediaReaders = map[string]RespBodyReader{
	TypeApplicationJson:   readAllRespWrapper(jsonRespBodyReader),
	TypeApplicationXml:    readAllRespWrapper(xmlRespBodyReader),
	TypeApplicationStream: streamRespBodyReader,
}

// readAllRespWrapper 从响应中读取全部数据的包装器
func readAllRespWrapper(readFunc func(resp *http.Response, data []byte, typ reflect.Type, out interface{}) error) RespBodyReader {
	return func(resp *http.Response, reader io.ReadCloser, typ reflect.Type, out interface{}) error {
		body, _ := ioutil.ReadAll(reader)
		return readFunc(resp, body, typ, out)
	}
}

// stringRespBodyReader 将数据读取到*string类型的out参数中
func stringRespBodyReader(resp *http.Response, data []byte, typ reflect.Type, out interface{}) error {
	*(out.(*string)) = string(data)
	return nil
}

// bytesRespBodyReader 将数据读取到*[]byte类型的out参数中
func bytesRespBodyReader(resp *http.Response, data []byte, typ reflect.Type, out interface{}) error {
	*(out.(*[]byte)) = data
	return nil
}

// jsonRespBodyReader 将json格式的数据的解析到out参数中
func jsonRespBodyReader(resp *http.Response, data []byte, typ reflect.Type, out interface{}) error {
	kind := typ.Elem().Kind()
	if kind == reflect.Struct || kind == reflect.Map {
		return json.Unmarshal(data, out)
	}
	return ErrRespOutType
}

// xmlRespBodyReader 将xml格式的数据的解析到out参数中
func xmlRespBodyReader(resp *http.Response, data []byte, typ reflect.Type, out interface{}) error {
	kind := typ.Elem().Kind()
	if kind == reflect.Struct || kind == reflect.Map {
		return xml.Unmarshal(data, out)
	}
	return ErrRespOutType
}

// responseReader 返回http response
func responseReader(resp *http.Response, reader io.ReadCloser, typ reflect.Type, out interface{}) error {
	output := out.(**http.Response)
	(*output) = resp
	return nil
}

// downloadRespBodyReader 将响应的数据保存到文件
func downloadRespBodyReader(resp *http.Response, reader io.ReadCloser, typ reflect.Type, out interface{}) error {
	var info *SaveInfo
	switch output := out.(type) {
	case FilePath:
		info = &SaveInfo{
			Path:     string(output),
			Override: true,
			Mode:     0666,
		}
	case SaveInfo:
		info = &output
	case *SaveInfo:
		info = output
	default:
		return ErrRespOutType
	}
	var file *os.File
	var err error
	if info.Override {
		file, err = os.OpenFile(info.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode)
	} else {
		file, err = os.OpenFile(info.Path, os.O_RDWR|os.O_CREATE|os.O_EXCL, info.Mode)
	}
	if err != nil {
		return err
	}
	file.Chmod(info.Mode)
	defer file.Close()
	written, err := io.Copy(file, resp.Body)
	if resp.ContentLength != written {
		return fmt.Errorf("save file size is %d, but content length is %d", written, resp.ContentLength)
	}
	return nil
}

// streamRespBodyReader 将响应的数据保存到文件
func streamRespBodyReader(resp *http.Response, reader io.ReadCloser, typ reflect.Type, out interface{}) error {
	kind := typ.Kind()
	if kind == reflect.String {
		return downloadRespBodyReader(resp, reader, typ, FilePath(out.(string)))
	}
	return ErrRespOutType
}
