package handlers

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Gopher0727/GoRepo/backend/middleware"
	"github.com/Gopher0727/GoRepo/backend/models"
	"github.com/Gopher0727/GoRepo/backend/security"
	"github.com/Gopher0727/GoRepo/backend/store"
)

// 简单内存用户存储（email -> encoded argon2 hash）
var (
	userStore   = map[string]string{}
	userStoreMu sync.RWMutex
)

type authRegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthRegister 处理用户注册
func AuthRegister(c *gin.Context) {
	var req authRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	userStoreMu.Lock()
	defer userStoreMu.Unlock()

	if _, ok := userStore[req.Email]; ok {
		c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
		return
	}

	hash, err := security.EncodeHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash failed"})
		return
	}
	userStore[req.Email] = hash

	// 写入数据库（若不存在）
	var u models.User
	if err := store.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u = models.User{Email: req.Email}
			if err2 := store.DB.Create(&u).Error; err2 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "db create user failed"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"email": req.Email, "id": u.ID})
}

// AuthLogin 处理用户登录
func AuthLogin(c *gin.Context) {
	var req authLoginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	userStoreMu.RLock()
	hash, ok := userStore[req.Email]
	userStoreMu.RUnlock()
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ok, err := security.VerifyHash(req.Password, hash)
	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 查询用户（若不存在则补建，避免旧数据缺失）
	var u models.User
	if err := store.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u = models.User{Email: req.Email}
			if err2 := store.DB.Create(&u).Error; err2 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "db create user failed"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
			return
		}
	}

	// 生成真实 JWT（sub 使用 email）
	tok, err := middleware.GenerateToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token gen failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tok, "email": req.Email, "id": u.ID})
}

// AuthRehashCheck 检查是否需要重新 hash（根据当前参数）
func AuthRehashCheck(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing email"})
		return
	}

	userStoreMu.RLock()
	hash, ok := userStore[email]
	userStoreMu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"needsRehash": security.NeedsRehash(hash)})
}
