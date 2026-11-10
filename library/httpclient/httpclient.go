package httpclient

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"template-go/utils/constant"
	"template-go/utils/tool/jsoniter"
)

type CommonResp[T any] struct {
	RequestID  string `json:"requestID"`
	Code       int    `json:"code"`    // 错误码
	Msg        string `json:"message"` // 错误信息
	Data       T      `json:"data"`    // 返回数据
	ServerTime int64  `json:"serverTime"`
}

func EmptyHeader() map[string]string {
	return map[string]string{}
}

func EmptyBody() []byte {
	return make([]byte, 0)
}

type HttpClient struct {
	client *http.Client
}

func CreateHttpClient() *HttpClient {
	transport := &http.Transport{
		MaxIdleConns:        90,
		MaxIdleConnsPerHost: 30,
		IdleConnTimeout:     30 * time.Second,
		DisableKeepAlives:   false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		//TLSClientConfig:       config.TLSClientConfig,
		ExpectContinueTimeout: time.Second,
	}

	return &HttpClient{&http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}}
}

func (cli *HttpClient) Do(ctx context.Context, module any, url, method string, headers map[string]string, body []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set(constant.HeaderContentType, constant.ContentTypeJSON)

	if headers == nil {
		headers = map[string]string{}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := cli.client.Do(req)
	if err != nil {
		return err
	}

	var respBody []byte
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.New(constant.GetCodeText(constant.CodeHTTPError))
	}

	err = jsoniter.Unmarshal(respBody, module)
	if err != nil {
		return err
	}

	return nil
}
