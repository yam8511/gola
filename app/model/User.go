package model

import (
	"gola/internal/database"
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

// TableName 資料表
func (m User) TableName() string {
	return TableUser
}

// Database 資料庫
func (m User) Database(master bool) database.Type {
	if master {
		return DBMaster
	}
	return DBSlave
}
