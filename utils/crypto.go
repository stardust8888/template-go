package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

// MD5 计算字符串的MD5值
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1 计算字符串的SHA1值
func SHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 计算字符串的SHA256值
func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA1 计算HMAC-SHA1签名
func HMACSHA1(message, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA256 计算HMAC-SHA256签名
func HMACSHA256(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Base64Encode Base64编码
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64Decode Base64解码
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

var ErrAESDecryptFailed = errors.New("decrypt failed")

// AESGCMEncrypt AES-GCM 加密（AEAD，密文自带完整性校验，被篡改会在解密时报错）
// 返回：nonce || ciphertext || tag
func AESGCMEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal 把密文和 tag 追加在 nonce 之后，nonce 作为前缀一起返回
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// AESGCMDecrypt AES-GCM 解密，密钥不对或密文被篡改都返回 ErrAESDecryptFailed
func AESGCMDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 至少要容得下 nonce 和 tag
	if len(ciphertext) < gcm.NonceSize()+gcm.Overhead() {
		return nil, ErrAESDecryptFailed
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrAESDecryptFailed
	}

	return plaintext, nil
}

// AESGCMEncryptHex 加密并编码成带 v2 前缀的 hex 字符串
func AESGCMEncryptHex(plaintext, key []byte) (string, error) {
	encrypt, err := AESGCMEncrypt(plaintext, key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypt), nil
}

// AESDecryptHex 解密 hex 密文：带 v2 前缀走 AES-GCM，不带前缀按历史 AES-CBC 兼容读取。
// 所有失败原因都归一成 ErrAESDecryptFailed。
func AESDecryptHex(ciphertextHex string, key []byte) ([]byte, error) {
	raw, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, ErrAESDecryptFailed
	}
	return AESGCMDecrypt(raw, key)
}

// GenerateKey 生成指定长度的随机密钥
func GenerateKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Hash 计算指定哈希算法的值
func Hash(data []byte, h hash.Hash) string {
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
