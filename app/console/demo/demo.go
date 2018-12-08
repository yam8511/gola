package demo

import (
	"gola/internal/bootstrap"
)

// Run 背景
func Run() error {
	bootstrap.WriteLog("INFO", "Demo Job")
	return nil
}
