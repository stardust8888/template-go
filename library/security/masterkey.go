// Package security 负责主密钥的加载与按用途派生。
//
// 主密钥来源优先级：
//  1. 环境变量 KMS_MASTER_KEY（生产环境唯一推荐方式，值为 32 字节密钥的 base64 编码）
//  2. 本地密钥文件 conf/.kms_master_key（环境变量未设置时自动生成并落盘，供开发环境使用）
//
// 环境变量已设置但格式非法时直接返回错误，不做任何降级——
// 静默退回默认密钥比启动失败危险得多，那会让人误以为数据已被正确加密。
package security

import (
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	// EnvMasterKey 主密钥环境变量名
	EnvMasterKey = "KMS_MASTER_KEY"

	// MasterKeySize 主密钥长度（AES-256）
	MasterKeySize = 32

	// keyFileDir/keyFileName 环境变量缺失时的本地落盘位置
	keyFileDir  = "conf"
	keyFileName = ".kms_master_key"

	// PurposeKMSCrypt KMS 加解密用途标签。
	// 这是密文格式的一部分：改动它会导致历史密文全部无法解密，必须保持稳定。
	PurposeKMSCrypt = "config-crypt"
)

var (
	initOnce   sync.Once
	initErr    error
	masterKey  []byte
	initSource string // "env" 或 "file"，仅用于日志

	deriveMu    sync.RWMutex
	deriveCache = map[string][]byte{}
)

// Init 加载主密钥，并发安全且只执行一次。
// 返回的错误应当被视为致命错误，调用方应拒绝启动。
func Init() error {
	initOnce.Do(func() {
		masterKey, initSource, initErr = load()
	})
	return initErr
}

// Source 返回主密钥的来源（env / file / generated），供启动日志区分环境。
func Source() string {
	return initSource
}

// Get 返回主密钥本身。多数场景应使用 DeriveKey 而非直接用主密钥。
func Get() ([]byte, error) {
	if err := Init(); err != nil {
		return nil, err
	}
	return masterKey, nil
}

// DeriveKey 用 HKDF-SHA256 从主密钥派生指定用途的子密钥（32 字节），结果缓存复用。
// purpose 是用途标签，不同用途得到互不相关的密钥，避免一把密钥跨场景复用。
func DeriveKey(purpose string) ([]byte, error) {
	if purpose == "" {
		return nil, errors.New("security: purpose 不能为空")
	}

	if err := Init(); err != nil {
		return nil, err
	}

	deriveMu.RLock()
	key, ok := deriveCache[purpose]
	deriveMu.RUnlock()
	if ok {
		return key, nil
	}

	derived, err := hkdf.Key(sha256.New, masterKey, nil, purpose, 32)
	if err != nil {
		return nil, err
	}

	deriveMu.Lock()
	deriveCache[purpose] = derived
	deriveMu.Unlock()

	return derived, nil
}

// load 按「环境变量 → 本地文件 → 自动生成」的顺序取主密钥。
func load() ([]byte, string, error) {
	if raw := strings.TrimSpace(os.Getenv(EnvMasterKey)); raw != "" {
		key, err := decodeKey(raw)
		if err != nil {
			return nil, "", fmt.Errorf("环境变量 %s 非法: %w", EnvMasterKey, err)
		}
		return key, "env", nil
	}

	// 环境变量没设置：先看有没有上次落盘的密钥
	for _, path := range keyFileCandidates() {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		key, err := decodeKey(strings.TrimSpace(string(raw)))
		if err != nil {
			return nil, "", fmt.Errorf("密钥文件 %s 非法: %w", path, err)
		}
		warnf("主密钥来自本地文件 %s（开发模式）。生产环境请改用环境变量 %s", path, EnvMasterKey)
		return key, "file", nil
	}

	// 都没有：生成一个并落盘，保证重启后仍能解开历史密文
	key := make([]byte, MasterKeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, "", fmt.Errorf("生成主密钥失败: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(key)
	path := keyFileWritePath()
	_ = os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, []byte(encoded), 0o600); err != nil {
		warnf("主密钥落盘失败（本次仅存在于内存，重启后历史密文将无法解密）: %v", err)
	} else {
		warnf("未设置环境变量 %s，已自动生成主密钥并写入 %s。生产环境请改用环境变量注入，并妥善备份该密钥", EnvMasterKey, path)
	}

	return key, "generated", nil
}

// decodeKey 校验密钥文本是 base64 编码的 32 字节值。
func decodeKey(raw string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("不是合法的 base64")
	}
	if len(key) != MasterKeySize {
		return nil, fmt.Errorf("长度必须为 %d 字节，当前 %d 字节", MasterKeySize, len(key))
	}
	return key, nil
}

// keyFileCandidates 密钥文件的读取位置，与配置文件搜索路径保持一致。
func keyFileCandidates() []string {
	return []string{
		filepath.Join(keyFileDir, keyFileName),
		filepath.Join("..", keyFileDir, keyFileName),
		filepath.Join("..", "..", keyFileDir, keyFileName),
		filepath.Join("..", "..", "..", keyFileDir, keyFileName),
	}
}

// keyFileWritePath 密钥文件的写入位置，优先放在 conf 目录下。
func keyFileWritePath() string {
	if stat, err := os.Stat(keyFileDir); err == nil && stat.IsDir() {
		return filepath.Join(keyFileDir, keyFileName)
	}
	return keyFileName
}

// warnf 启动告警直接写 stderr——此时日志系统未必已初始化，且这条信息不能被日志级别吞掉。
func warnf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[security] "+format+"\n", args...)
}
