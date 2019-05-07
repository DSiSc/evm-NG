package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/DSiSc/evm-NG/system"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

//Cos response error
type RespError struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	Resource  string `xml:"Resource"`
	RequestId string `xml:"RequestId"`
	TraceId   string `xml:"TraceId"`
}

// ObjectMeta object meta info
type ObjectMeta struct {
	ETag          string
	VersionId     string
	EncryptionAlg string
}

// tencent cos call routes
var CosRoutes = map[string]*system.SysCallFunc{
	"GetObject":         system.NewSysCallFunc(GetObject, "rawurl,name"),
	"PutObject":         system.NewSysCallFunc(PutObject, "rawurl,name,objBytes"),
	"PutObjectFromFile": system.NewSysCallFunc(PutObjectFromFile, "rawurl,name,objBytes"),
}

// GetObject download an object from the cloud server
func GetObject(rawurl, name string) ([]byte, error) {
	client, err := buildCosClient(rawurl)
	if err != nil {
		return nil, err
	}

	resp, err := client.Object.Get(context.Background(), name, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err = checkResponse(resp); err != nil {
		return nil, err
	}

	var respBytes []byte
	err = parseResp(resp.Body, bytesType, &respBytes)
	return respBytes, err
}

// PutObject upload an object to cloud server
func PutObject(rawurl, name string, objBytes []byte) (*ObjectMeta, error) {
	client, err := buildCosClient(rawurl)
	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(objBytes)
	resp, err := client.Object.Put(context.Background(), name, br, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err = checkResponse(resp); err != nil {
		return nil, err
	}

	objMeta := getObjectMeta(resp.Header)
	return objMeta, err
}

// PutObjectFromFile upload local file to cloud server
func PutObjectFromFile(rawurl, name, filePath string) (*ObjectMeta, error) {
	client, err := buildCosClient(rawurl)
	if err != nil {
		return nil, err
	}

	resp, err := client.Object.PutFromFile(context.Background(), name, filePath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err = checkResponse(resp); err != nil {
		return nil, err
	}

	objMeta := getObjectMeta(resp.Header)
	return objMeta, err
}

// build cos client with specified url
func buildCosClient(rawurl string) (*cos.Client, error) {
	u, e := url.Parse(rawurl)
	if e != nil {
		return nil, e
	}
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{})
	return c, nil
}

// extract object meta info from response header
func getObjectMeta(header http.Header) *ObjectMeta {
	return &ObjectMeta{
		ETag:          header.Get("ETag"),
		VersionId:     header.Get("x-cos-version-id"),
		EncryptionAlg: header.Get("x-cos-server-side-encryption"),
	}
}

// check response status
func checkResponse(resp *cos.Response) (err error) {
	if 200 != resp.StatusCode {
		var respError RespError
		if err1 := parseResp(resp.Body, xmlType, &respError); err1 != nil {
			return err1
		} else {
			return errors.Errorf("response error, Code: %s, Message: %s, Resource: %s, RequestId: %s, TraceId: %s", respError.Code, respError.Message, respError.Resource, respError.RequestId, respError.TraceId)
		}
	} else {
		return nil
	}
}

const (
	xmlType   = "xml"
	jsonType  = "json"
	bytesType = "bytes"
)

// parse response
func parseResp(resp io.Reader, parseType string, v interface{}) error {
	respBytes, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}
	switch parseType {
	case xmlType:
		return xml.Unmarshal(respBytes, v)
	case jsonType:
		return json.Unmarshal(respBytes, v)
	case bytesType:
		srcV, dstV := reflect.ValueOf(v), reflect.ValueOf(&respBytes)
		srcV.Elem().Set(dstV.Elem())
		return nil
	default:
		return errors.New("unknown parse type")
	}
}
