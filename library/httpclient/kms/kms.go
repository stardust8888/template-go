package kms

import (
	"context"
	"errors"
	"net/http"

	"template-go/library/httpclient"
	"template-go/utils/constant"
	"template-go/utils/tool/jsoniter"
)

var (
	client *Client
)

type Client struct {
	Endpoint   string
	HttpClient *httpclient.HttpClient
}

func InitClient() {
	client = &Client{
		Endpoint:   "http://127.0.0.1:12006",
		HttpClient: httpclient.CreateHttpClient(),
	}
}

func GetClient() *Client {
	return client
}

type EncryptReq struct {
	Plaintext string `json:"plaintext"`
}

type EncryptResp struct {
	Ciphertext string `json:"ciphertext"`
}

// Encrypt 调用 kms 加密，返回 hex 密文
func (c *Client) Encrypt(ctx context.Context, plaintext string) (string, error) {
	var (
		resp = httpclient.CommonResp[EncryptResp]{}
		err  error
	)
	if plaintext == "" {
		return "", errors.New("kms: 参数不能为空")
	}

	body, err := jsoniter.Marshal(EncryptReq{Plaintext: plaintext})
	if err != nil {
		return "", err
	}
	fullURL := c.Endpoint + constant.BaseKMSUrlV1 + "/crypt"

	if err = c.HttpClient.Do(ctx, &resp, fullURL, http.MethodPost, nil, body); err != nil {
		return "", err
	}

	if resp.Code != constant.CodeSuccess {
		return "", errors.New("encrypt error")
	}

	return resp.Data.Ciphertext, nil
}

type DecryptReq struct {
	Ciphertext string `json:"ciphertext"`
}

type DecryptResp struct {
	Plaintext string `json:"plaintext"`
}

// Decrypt 调用 kms 解密，传入 hex 密文，返回明文
func (c *Client) Decrypt(ctx context.Context, ciphertext string) (string, error) {
	var (
		resp = httpclient.CommonResp[DecryptResp]{}
		err  error
	)
	if ciphertext == "" {
		return "", errors.New("kms: 参数不能为空")
	}

	body, err := jsoniter.Marshal(DecryptReq{Ciphertext: ciphertext})
	if err != nil {
		return "", err
	}
	fullURL := c.Endpoint + constant.BaseKMSUrlV1 + "/decrypt"

	if err = c.HttpClient.Do(ctx, &resp, fullURL, http.MethodPost, nil, body); err != nil {
		return "", err
	}

	if resp.Code != constant.CodeSuccess {
		return "", errors.New("decrypt error")
	}

	return resp.Data.Plaintext, nil
}
