package demo

import (
	"gola/internal/logger"
	"time"

	"github.com/spf13/cobra"
)

// Run 背景
func Run(cmd *cobra.Command, args []string) error {
	logger.Info("Demo Job")
	time.Sleep(time.Second * 5)
	return nil
}
