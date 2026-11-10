package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"template-go/utils/tool/jsoniter"
	"template-go/utils/tool/spew"
)

func ToRawMessage(i any) json.RawMessage {
	b := ToJsonByte(i)
	return b
}

func ToJsonByte(i any) []byte {
	b, _ := jsoniter.Marshal(i)
	return b
}

func ToJson(i any) string {
	b, _ := jsoniter.Marshal(i)
	return string(b)
}

func ToString(i any) string {
	return fmt.Sprint(i)
}

func Spew(i any) {
	spew.Dump(i)
}

const (
	AlphaNumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	NumberCharset       = "012356789" // 有意不含 4
)

// RandString 生成指定长度的随机字母数字串（crypto/rand，可用于密钥、token、盐）
func RandString(length int) string {
	return RandStringWithCharset(length, AlphaNumericCharset)
}

// RandNumberString 生成指定长度的随机数字串（crypto/rand）
func RandNumberString(length int) string {
	return RandStringWithCharset(length, NumberCharset)
}

// RandStringWithCharset 使用指定字符集生成指定长度的随机字符串。
// 基于 crypto/rand，并用拒绝采样保证每个字符出现概率相同。
func RandStringWithCharset(length int, charset string) string {
	n := len(charset)
	if length <= 0 || n == 0 {
		return ""
	}

	out := make([]byte, 0, length)

	// 字符集超过 256 个字节时，单个随机字节不够用，逐个取 [0, n) 的随机数
	if n > 256 {
		max := big.NewInt(int64(n))
		for len(out) < length {
			idx, err := rand.Int(rand.Reader, max)
			if err != nil {
				panic("crypto/rand: " + err.Error())
			}
			out = append(out, charset[idx.Int64()])
		}
		return string(out)
	}

	// 只接受 < limit 的字节，limit 是 n 的整数倍，避免取模带来的偏差
	limit := 256 - 256%n
	buf := make([]byte, length+length/2+8)
	for len(out) < length {
		if _, err := rand.Read(buf); err != nil {
			panic("crypto/rand: " + err.Error())
		}
		for _, b := range buf {
			if int(b) >= limit {
				continue
			}
			out = append(out, charset[int(b)%n])
			if len(out) == length {
				break
			}
		}
	}
	return string(out)
}
