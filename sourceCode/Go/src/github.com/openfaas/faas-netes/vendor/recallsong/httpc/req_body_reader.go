package httpc

// Author: recallsong
// Email: songruiguo@qq.com

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

// ReqBodyReader 根据类型将数据转换为reader
type ReqBodyReader func(req *http.Request, typ reflect.Type, data interface{}) (io.ReadCloser, error)

var (
	// ErrReqBodyType 请求的Body数据类型错误
	ErrReqBodyType = errors.New("invalid request body type")
)

// GlobalReqBodyTypeReaders 默认的根据类型获取RequestBody的Reader映射
var GlobalReqBodyTypeReaders = map[reflect.Type]ReqBodyReader{
	reflect.TypeOf(""):              stringReqBodyReader,
	reflect.TypeOf([]byte(nil)):     bytesReqBodyReader,
	reflect.TypeOf(url.Values(nil)): urlValuesReqBodyReader,
}

// GlobalReqBodyMediaReaders 默认的根据MediaType获取RequestBody的Reader映射
var GlobalReqBodyMediaReaders = map[string]ReqBodyReader{
	TypeApplicationJson:   jsonReqBodyReader,
	TypeApplicationXml:    xmlReqBodyReader,
	TypeApplicationForm:   formValuesReqBodyReader,
	TypeMultipartFormData: multipartReqBodyReader,
}

// stringReqBodyReader 将string类型的data转换为reader
func stringReqBodyReader(req *http.Request, typ reflect.Type, data interface{}) (io.ReadCloser, error) {
	body := data.(string)
	req.ContentLength = int64(len(body))
	return ioutil.NopCloser(strings.NewReader(body)), nil
}

// bytesReqBodyReader 将[]byte类型的data转换为reader
func bytesReqBodyReader(req *http.Request, typ reflect.Type, data interface{}) (io.ReadCloser, error) {
	body := data.([]byte)
	req.ContentLength = int64(len(body))
	return ioutil.NopCloser(bytes.NewReader(body)), nil
}

// urlValuesReqBodyReader 将url.Values类型的data转换为reader
func urlValuesReqBodyReader(req *http.Request, typ reflect.Type, data interface{}) (io.ReadCloser, error) {
	body := (data.(url.Values)).Encode()
	req.ContentLength = int64(len(body))
	return ioutil.NopCloser(strings.NewReader(body)), nil
}

// jsonReqBodyReader 将data转换为json格式，并返回对应的reader
func jsonReqBodyReader(req *http.Request, typ reflect.Type, data interface{}) (io.ReadCloser, error) {
	contentJSON, err := json.Marshal(data)
	if err != nil {
		return nil, ErrReqBodyType
	}
	req.ContentLength = int64(len(contentJSON))
	return ioutil.NopCloser(bytes.NewReader(contentJSON)), nil
}

// xmlReqBodyReader 将data转换为xml格式，并返回对应的reader
func xmlReqBodyReader(req *http.Request, typ reflect.Type, data interface{}) (io.ReadCloser, error) {
	contentXML, err := xml.Marshal(data)
	if err != nil {
		return nil, ErrReqBodyType
	}
	req.ContentLength = int64(len(contentXML))
	return ioutil.NopCloser(bytes.NewReader(contentXML)), nil
}

// formValuesReqBodyReader 将data转换为form参数格式，并返回对应的reader
func formValuesReqBodyReader(req *http.Request, typ reflect.Type, data interface{}) (io.ReadCloser, error) {
	kind := typ.Kind()
	if kind == reflect.Map {
		switch v := data.(type) {
		case map[string]string:
			params := url.Values{}
			for key, val := range v {
				params.Add(key, val)
			}
			body := params.Encode()
			req.ContentLength = int64(len(body))
			return ioutil.NopCloser(strings.NewReader(body)), nil
		case map[string][]string:
			body := url.Values(v).Encode()
			req.ContentLength = int64(len(body))
			return ioutil.NopCloser(strings.NewReader(body)), nil
		case map[string]interface{}:
			params := url.Values{}
			for key, val := range v {
				if val != nil {
					params.Add(key, fmt.Sprint(val))
				}
			}
			body := params.Encode()
			req.ContentLength = int64(len(body))
			return ioutil.NopCloser(strings.NewReader(body)), nil
		default:
			return nil, ErrReqBodyType
		}
	} else if kind == reflect.Struct || (kind == reflect.Ptr && typ.Elem().Kind() == reflect.Struct) {
		params := url.Values{}
		value := reflect.ValueOf(data)
		if kind == reflect.Ptr {
			value = value.Elem()
			typ = typ.Elem()
		}
		num := typ.NumField()
		for i := 0; i < num; i++ {
			t := typ.Field(i)
			v := value.Field(i)
			tagVal := t.Tag.Get("form")
			if v.CanInterface() && tagVal != "" && tagVal != "-" {
				val := v.Interface()
				if val != nil {
					params.Add(tagVal, fmt.Sprint(val))
				}
			}
		}
		body := params.Encode()
		req.ContentLength = int64(len(body))
		return ioutil.NopCloser(strings.NewReader(body)), nil
	}
	return nil, ErrReqBodyType
}

// getFileShortName 获取简短的文件名
func getFileShortName(file *os.File) string {
	longName := file.Name()
	return longName[strings.LastIndex(longName, "/")+1:]
}

// writeToMultipartFormData 将数据写入writer中
func writeToMultipartFormData(key string, val interface{}, fileName string, writer *multipart.Writer, closeList *[]*os.File) error {
	switch realVal := val.(type) {
	case string:
		err := writer.WriteField(key, realVal)
		if err != nil {
			return err
		}
	case []string:
		for _, val := range realVal {
			err := writer.WriteField(key, val)
			if err != nil {
				return err
			}
		}
	case *os.File:
		if realVal == nil {
			return nil
		}
		if fileName == "" {
			fileName = getFileShortName(realVal)
		}
		part, err := writer.CreateFormFile(key, fileName)
		if err != nil {
			return err
		}
		_, err = io.Copy(part, realVal)
		if err != nil {
			return err
		}
	case FilePath:
		if realVal == "" {
			return nil
		}
		path := string(realVal)
		idx := strings.LastIndex(path, ":")
		if idx >= 0 {
			fileName = path[idx+1:]
			path = path[0:idx]
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		*closeList = append(*closeList, file)
		if fileName == "" {
			fileName = getFileShortName(file)
		}
		part, err := writer.CreateFormFile(key, fileName)
		if err != nil {
			return err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}
	default:
		if val != nil {
			err := writer.WriteField(key, fmt.Sprint(val))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// fileReqBodyReader 将data转换为multipart格式，并返回对应的reader
func multipartReqBodyReader(req *http.Request, typ reflect.Type, data interface{}) (reader io.ReadCloser, err error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	closeList := make([]*os.File, 0)
	go func() (err error) {
		defer func() {
			var errs string
			werr := writer.Close()
			if werr != nil {
				errs += fmt.Sprintf("multipart.Writer.Close: %s\n", werr.Error())
			}
			for _, item := range closeList {
				err := item.Close()
				if err != nil {
					errs += fmt.Sprintf("close %s: %s\n", item.Name(), err.Error())
				}
			}
			if err != nil {
				errs += fmt.Sprintf("process: %s\n", err.Error())
			}
			if errs != "" {
				errs = errs[0 : len(errs)-1]
				pw.CloseWithError(errors.New(errs))
			} else {
				pw.Close()
			}
		}()
		kind := typ.Kind()
		if kind == reflect.Map {
			switch v := data.(type) {
			case map[string]string:
				for key, val := range v {
					err = writer.WriteField(key, val)
					if err != nil {
						return err
					}
				}
			case map[string][]string:
				for key, vals := range v {
					for _, val := range vals {
						err = writer.WriteField(key, val)
						if err != nil {
							return err
						}
					}
				}
			case map[string]interface{}:
				for key, val := range v {
					err = writeToMultipartFormData(key, val, "", writer, &closeList)
					if err != nil {
						return err
					}
				}
			case map[string]FilePath:
				for key, val := range v {
					err = writeToMultipartFormData(key, val, "", writer, &closeList)
					if err != nil {
						return err
					}
				}
			case map[string]*os.File:
				for key, val := range v {
					err = writeToMultipartFormData(key, val, "", writer, &closeList)
					if err != nil {
						return err
					}
				}
			default:
				return ErrReqBodyType
			}
		} else if kind == reflect.Struct || (kind == reflect.Ptr && typ.Elem().Kind() == reflect.Struct) {
			value := reflect.ValueOf(data)
			if kind == reflect.Ptr {
				value = value.Elem()
				typ = typ.Elem()
			}
			num := typ.NumField()
			for i := 0; i < num; i++ {
				t := typ.Field(i)
				v := value.Field(i)
				tagVal := t.Tag.Get("form")
				fileName := t.Tag.Get("file")
				if v.CanInterface() && tagVal != "" && tagVal != "-" {
					if fileName == "" || fileName == "-" {
						fileName = ""
					}
					err = writeToMultipartFormData(tagVal, v.Interface(), fileName, writer, &closeList)
					if err != nil {
						return err
					}
				}
			}
		} else {
			return ErrReqBodyType
		}
		return
	}()
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return pr, nil
}
