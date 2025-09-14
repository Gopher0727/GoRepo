package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Gopher0727/GoRepo/backend/models"
	"github.com/Gopher0727/GoRepo/backend/store"
)

// helper: 获取当前用户（通过中间件放入的 sub=email）
func getCurrentUser(c *gin.Context) (*models.User, bool) {
	sub, ok := c.Get("sub")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no subject"})
		return nil, false
	}

	email, _ := sub.(string)
	var u models.User
	if err := store.DB.Where("email = ?", email).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return nil, false
	}
	return &u, true
}

// EntryList 列出当前用户所有条目
func EntryList(c *gin.Context) {
	u, ok := getCurrentUser(c)
	if !ok {
		return
	}

	var list []models.Entry
	if err := store.DB.Where("user_id = ?", u.ID).Order("id desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list})
}

// EntryCreate 创建条目
func EntryCreate(c *gin.Context) {
	u, ok := getCurrentUser(c)
	if !ok {
		return
	}

	var req models.Entry
	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	entry := models.Entry{
		UserID:   u.ID,
		Title:    req.Title,
		Username: req.Username,
		URL:      req.URL,
	}
	if err := store.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db create failed"})
		return
	}
	c.JSON(http.StatusCreated, entry)
}

// EntryDetail 获取条目详情
func EntryDetail(c *gin.Context) {
	u, ok := getCurrentUser(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	var e models.Entry
	if err := store.DB.Where("id = ? AND user_id = ?", id, u.ID).First(&e).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, e)
}

// EntryUpdate 更新条目
func EntryUpdate(c *gin.Context) {
	u, ok := getCurrentUser(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	var req models.Entry
	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	var e models.Entry
	if err := store.DB.Where("id = ? AND user_id = ?", id, u.ID).First(&e).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	e.Title = req.Title
	e.Username = req.Username
	e.URL = req.URL
	if err := store.DB.Save(&e).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db save failed"})
		return
	}
	c.JSON(http.StatusOK, e)
}

// EntryDelete 删除条目
func EntryDelete(c *gin.Context) {
	u, ok := getCurrentUser(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	if err := store.DB.Where("id = ? AND user_id = ?", id, u.ID).Delete(&models.Entry{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}
