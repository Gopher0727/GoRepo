package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Gopher0727/GoRepo/backend/middleware"
	"github.com/Gopher0727/GoRepo/backend/models"
	"github.com/Gopher0727/GoRepo/backend/security"
	"github.com/Gopher0727/GoRepo/backend/store"
)

type authRegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type authLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthRegister 处理用户注册
func AuthRegister(c *gin.Context) {
	var req authRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	var existed models.User
	if err := store.DB.Where("email = ?", req.Email).First(&existed).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	hash, err := security.EncodeHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash failed"})
		return
	}

	// 写入数据库（若不存在）
	u := models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: hash,
	}

	if err := store.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db create user failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
	})
}

// AuthLogin 处理用户登录
func AuthLogin(c *gin.Context) {
	var req authLoginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	var u models.User
	if err := store.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	ok, err := security.VerifyHash(req.Password, u.Password)
	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	tok, err := middleware.GenerateToken(u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token gen failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tok,
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
	})
}

// AuthRehashCheck 检查是否需要重新 hash（根据当前参数）
func AuthRehashCheck(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing email"})
		return
	}

	var u models.User
	if err := store.DB.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"needsRehash": security.NeedsRehash(u.Password)})
}
