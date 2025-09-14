package models

import "time"

// Entry 密码条目
type Entry struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"-"`
	Title     string    `gorm:"type:varchar(150);not null;index:idx_user_title,priority:2" json:"title"`
	Username  string    `gorm:"type:varchar(150);not null;index:idx_user_title,priority:3" json:"username"`
	URL       string    `gorm:"type:varchar(300);not null;default:''" json:"url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
