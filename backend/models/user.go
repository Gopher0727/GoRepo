package models

import "time"

// User 用户模型
type User struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string `gorm:"type:varchar(191);uniqueIndex;not null" json:"email"`
	Name      string `gorm:"size:64;not null;default:''"`
	Password  string `gorm:"size:191;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
