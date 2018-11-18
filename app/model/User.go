package model

import (
	"time"
)

// User 使用者
type User struct {
	ID        int64
	Username  string // 帳號
	Name      string // 名稱
	Password  string // 密碼
	Enable    bool   // 啟用狀態
	CreatedAt time.Time
	UpdatedAt time.Time
}
