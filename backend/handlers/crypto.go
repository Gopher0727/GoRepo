package handlers

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"

	"github.com/Gopher0727/GoRepo/backend/security"
)

type encryptReq struct {
	Key       string `json:"key"`       // base64
	Plaintext string `json:"plaintext"` // 原文（直接字符串）
}

type encryptResp struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
}

type decryptReq struct {
	Key        string `json:"key"`
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
}

// EncryptData Gin 处理函数
func EncryptData(c *gin.Context) {
	var req encryptReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Key == "" {
		c.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	keyBytes, err := base64.StdEncoding.DecodeString(req.Key)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad key b64"})
		return
	}

	ct, nonce, err := security.Encrypt(keyBytes, []byte(req.Plaintext))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, encryptResp{
		Ciphertext: base64.StdEncoding.EncodeToString(ct),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
	})
}

// DecryptData Gin 处理函数
func DecryptData(c *gin.Context) {
	var req decryptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	keyBytes, err := base64.StdEncoding.DecodeString(req.Key)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad key b64"})
		return
	}

	ct, err := base64.StdEncoding.DecodeString(req.Ciphertext)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad ciphertext b64"})
		return
	}

	nonce, err := base64.StdEncoding.DecodeString(req.Nonce)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad nonce b64"})
		return
	}

	pt, err := security.Decrypt(keyBytes, ct, nonce)
	if err != nil {
		c.JSON(400, gin.H{"error": "decrypt failed"})
		return
	}
	c.JSON(200, gin.H{"plaintext": string(pt)})
}
