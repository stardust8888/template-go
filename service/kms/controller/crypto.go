package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"template-go/library/ginx"
	"template-go/library/security"
	"template-go/logs"
	"template-go/utils"
)

// cryptKey 取 KMS 专用的加密子密钥。密钥来自环境变量 KMS_MASTER_KEY 或本地密钥文件，
// 由 library/security 在启动时加载，缺失时自动生成并落盘。
func cryptKey() ([]byte, error) {
	return security.DeriveKey(security.PurposeKMSCrypt)
}

type EncryptReq struct {
	Plaintext string `json:"plaintext" binding:"required"`
}

type EncryptResp struct {
	Ciphertext string `json:"ciphertext"`
}

func Encrypt(c *gin.Context) {
	ctx := c.Request.Context()
	req := EncryptReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.ResponseParamError(c)
		return
	}

	key, err := cryptKey()
	if err != nil {
		logs.Error(ctx, "加载加密密钥失败", zap.Error(err))
		ginx.ResponseServerError(c, err)
		return
	}

	ciphertext, err := utils.AESGCMEncryptHex([]byte(req.Plaintext), key)
	if err != nil {
		logs.Error(ctx, "加密失败", zap.Error(err))
		ginx.ResponseServerError(c, err)
		return
	}

	resp := EncryptResp{
		Ciphertext: ciphertext,
	}
	ginx.ResponseSuccess(c, resp)
}

type DecryptReq struct {
	Ciphertext string `json:"ciphertext" binding:"required"`
}

type DecryptResp struct {
	Plaintext string `json:"plaintext"`
}

func Decrypt(c *gin.Context) {
	ctx := c.Request.Context()
	req := DecryptReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.ResponseParamError(c)
		return
	}

	// 解密失败（hex 格式不对 / 密钥不匹配 / 密文被篡改）一律返回 CodeParamError，
	// 不区分具体原因，避免调用方利用返回码差异做 padding oracle 攻击。
	key, err := cryptKey()
	if err != nil {
		logs.Error(ctx, "加载解密密钥失败", zap.Error(err))
		ginx.ResponseServerError(c, err)
		return
	}

	plaintext, err := utils.AESDecryptHex(req.Ciphertext, key)
	if err != nil {
		logs.Warn(ctx, "解密失败", zap.Error(err))
		ginx.ResponseParamError(c)
		return
	}

	resp := DecryptResp{
		Plaintext: string(plaintext),
	}
	ginx.ResponseSuccess(c, resp)
}
