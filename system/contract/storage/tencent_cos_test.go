package storage

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/DSiSc/monkey"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

const (
	rawUrl  = "https://blockchain-contract-1259100639.cos.ap-beijing.myqcloud.com"
	objName = "hello.txt"
)

func mockTencentCosContract() *TencentCosContract {
	return &TencentCosContract{
		rw: bytes.NewBuffer(make([]byte, 1024)),
	}
}

func mockClient() *cos.Client {
	client, _ := buildCosClient(rawUrl)
	return client
}

func mockSuccessResponse() *cos.Response {
	return &cos.Response{
		Response: &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("Hello")),
		},
	}
}

var respError = RespError{
	Code:      "400",
	Message:   "Error Addr",
	Resource:  "http://hello.test/testobj",
	RequestId: "1234",
	TraceId:   "1234",
}
var mockRespError = errors.Errorf("response error, Code: %s, Message: %s, Resource: %s, RequestId: %s, TraceId: %s", respError.Code, respError.Message, respError.Resource, respError.RequestId, respError.TraceId)

func mockErrorResponse() *cos.Response {
	errBytes, _ := xml.Marshal(respError)
	return &cos.Response{
		Response: &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader(errBytes)),
		},
	}
}

func TestGetObject(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	tencentCos := mockTencentCosContract()
	buf := bytes.NewBufferString("")
	tencentCos.rw = buf
	client := mockClient()
	monkey.Patch(cos.NewClient, func(uri *cos.BaseURL, httpClient *http.Client) *cos.Client {
		return client
	})
	response := mockSuccessResponse()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.Object), "Get", func(obj *cos.ObjectService, ctx context.Context, name string, opt *cos.ObjectGetOptions, id ...string) (*cos.Response, error) {
		return response, nil
	})
	err := tencentCos.GetObject(rawUrl, objName)
	assert.Nil(err)
	obj, _ := ioutil.ReadAll(tencentCos.rw)
	assert.Equal([]byte("Hello"), obj)

	response = mockErrorResponse()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.Object), "Get", func(obj *cos.ObjectService, ctx context.Context, name string, opt *cos.ObjectGetOptions, id ...string) (*cos.Response, error) {
		return response, nil
	})
	err = tencentCos.GetObject(rawUrl, objName)
	assert.EqualError(err, fmt.Sprintf("response error, Code: %s, Message: %s, Resource: %s, RequestId: %s, TraceId: %s", respError.Code, respError.Message, respError.Resource, respError.RequestId, respError.TraceId))
}

func TestPutObject(t *testing.T) {
	defer monkey.UnpatchAll()
	objBytes := []byte{'h', 'e', 'l', 'l', 'o'}
	objMeta := &ObjectMeta{
		ETag: "b1946ac92492d2347c6235b4d2611184",
	}
	assert := assert.New(t)
	tencentCos := mockTencentCosContract()
	tencentCos.rw.Write(objBytes)
	client := mockClient()
	monkey.Patch(cos.NewClient, func(uri *cos.BaseURL, httpClient *http.Client) *cos.Client {
		return client
	})

	response := mockSuccessResponse()
	response.Header = http.Header{}
	response.Header.Set("ETag", "b1946ac92492d2347c6235b4d2611184")
	monkey.PatchInstanceMethod(reflect.TypeOf(client.Object), "Put", func(obj *cos.ObjectService, ctx context.Context, name string, r io.Reader, opt *cos.ObjectPutOptions) (*cos.Response, error) {
		return response, nil
	})
	objMeta1, err := tencentCos.PutObject(rawUrl, objName)
	assert.Nil(err)
	assert.Equal(objMeta.ETag, objMeta1.ETag)

	response = mockErrorResponse()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.Object), "Put", func(obj *cos.ObjectService, ctx context.Context, name string, r io.Reader, opt *cos.ObjectPutOptions) (*cos.Response, error) {
		return response, nil
	})
	_, err = tencentCos.PutObject(rawUrl, objName)
	assert.EqualError(err, fmt.Sprintf("response error, Code: %s, Message: %s, Resource: %s, RequestId: %s, TraceId: %s", respError.Code, respError.Message, respError.Resource, respError.RequestId, respError.TraceId))
}
