package demo

import (
	"gola/internal/logger"
	"time"
)

// Run 背景
func Run() error {
	logger.Info("Demo Job")
	time.Sleep(time.Second * 5)
	return nil
}
