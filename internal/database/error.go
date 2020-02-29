package database

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// 連接池相關錯誤
var (
	ErrPoolHasNoConf = errors.New("連線池尚未設定")
	ErrPoolHasClosed = errors.New("連線池已經關閉")
	ErrPoolTimeout   = errors.New("逾期5秒，請稍後再取DB連線")
)

// 回傳尚未設定的錯誤，並且印出Log
func newNoConfError(name string) error {
	log.Printf("[%s] 連線池尚未設定", name)
	return ErrPoolHasNoConf
}

// IsPoolClosed 連線池是否關閉
func IsPoolClosed(err error) bool {
	if err == nil {
		return false
	}

	switch err {
	case ErrPoolHasClosed,
		sql.ErrConnDone,
		mysql.ErrInvalidConn:
		return true
	}

	errText := err.Error()

	if strings.Contains(errText, "connect: connection refused") {
		return true
	}

	switch errText {
	case "redis: client is closed",
		"sql: database is closed":
		return true
	}

	return false
}

// IsPoolTimeout 連線池是否逾時
func IsPoolTimeout(err error) bool {
	if err == nil {
		return false
	}

	switch err {
	case ErrPoolTimeout,
		mysql.ErrBusyBuffer,
		driver.ErrBadConn:
		return true
	}

	errText := err.Error()
	if errText == "redis: connection pool timeout" {
		return true
	}

	return false
}
