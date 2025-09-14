package models

// User 用户模型
type User struct {
	ID    uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Email string `gorm:"type:varchar(191);uniqueIndex;not null" json:"email"`
}
