package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 参数（后续可在配置里升级）
const (
	argonTime    uint32 = 1         // 迭代次数（可调大提升安全性）
	argonMemory  uint32 = 64 * 1024 // KiB -> 64MB
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32 // 生成 32 字节 key (可作 AES-256 密钥)
	argonSaltLen uint32 = 16
)

// Derive 从 password + salt 派生 key；若 salt == nil 自动生成随机盐
func Derive(password string, salt []byte) (key, outSalt []byte, err error) {
	if password == "" {
		return nil, nil, errors.New("empty password")
	}
	if salt == nil {
		salt = make([]byte, argonSaltLen)
		if _, err = rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}
	key = argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return key, salt, nil
}

// EncodeHash 生成可持久化的 Argon2id 字符串（包含参数/盐/哈希）
// 规范格式: $argon2id$v=19$m=65536,t=1,p=4$<saltBase64>$<hashBase64>
func EncodeHash(password string) (string, error) {
	key, salt, err := Derive(password, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

// VerifyHash 验证明文密码与编码串是否匹配
func VerifyHash(password, encoded string) (bool, error) {
	// 简单解析（不做完整语法验证）
	// 预期分段："", "argon2id", "v=19", "m=...,t=...,p=...", <saltB64>, <hashB64>
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid argon2 encoded format")
	}
	paramPart := parts[3]
	var m, t uint32
	var p uint8
	_, err := fmt.Sscanf(paramPart, "m=%d,t=%d,p=%d", &m, &t, &p)
	if err != nil {
		return false, errors.New("invalid argon2 params")
	}
	saltB, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	// 重新派生
	cand := argon2.IDKey([]byte(password), saltB, t, m, p, uint32(len(expectedHash)))
	if len(cand) != len(expectedHash) {
		return false, nil
	}
	// 常量时间比较
	var diff byte
	for i := range cand {
		diff |= cand[i] ^ expectedHash[i]
	}
	return diff == 0, nil
}

// NeedsRehash 判断现有编码是否需要根据当前参数重新计算
func NeedsRehash(encoded string) bool {
	// 简单解析（不做完整语法验证）
	// 预期分段："", "argon2id", "v=19", "m=...,t=...,p=...", <saltB64>, <hashB64>
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return true
	}
	if parts[2] != "v=19" { // 目前仅支持 v=19
		return true
	}

	// 解析参数
	var m, t uint32
	var p uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return true
	}
	return m != argonMemory || t != argonTime || p != argonThreads
}
