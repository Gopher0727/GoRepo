package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Encrypt 使用 AES-GCM 加密明文，返回 (ciphertext, nonce)
// key 长度需为 16/24/32 字节（AES-128/192/256）
func Encrypt(key, plaintext []byte) (ct, nonce []byte, err error) {
	// 校验长度
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, nil, errors.New("invalid aes key length")
	}

	// 创建分组密码实例：得到一个实现 Block 接口的 AES
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// 包装成 AEAD：获得一个支持认证加密的 GCM 对象
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	// 生成随机 nonce
	// - rand.Reader 是一个全局的、并发安全的密码用强随机数生成器
	nonce = make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	// 加密
	// - 用计数器模式（CTR 变体）对明文加密
	// - 计算 GHASH 认证标签（完整性与真实性）
	ct = gcm.Seal(nil, nonce, plaintext, nil)

	// 调用方需要保存 nonce，以便解密时使用
	return ct, nonce, nil
}

// Decrypt 使用 AES-GCM 解密，返回明文
func Decrypt(key, ciphertext, nonce []byte) ([]byte, error) {
	// 校验长度
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid aes key length")
	}

	// 创建分组密码实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 包装成 AEAD
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 解密
	// - 校验 GHASH 认证标签
	// - 用 CTR 变体对密文解密
	if len(nonce) != gcm.NonceSize() {
		return nil, errors.New("invalid nonce size")
	}

	// Open 会在认证失败时返回错误
	pt, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decrypt failed")
	}

	return pt, nil
}
