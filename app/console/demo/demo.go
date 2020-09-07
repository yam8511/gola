package demo

import (
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// Run 背景
func Run(cmd *cobra.Command, args []string) error {
	wait := make(chan byte)
	go func() {
		select {
		case <-wait:
		case <-bootstrap.GracefulDown():
			logger.Warn(cmd.Name() + "接收到中斷訊號，再收到第二次訊號，將強制結束程序！")
			select {
			case <-wait:
			case <-bootstrap.WaitOnceSignal():
				os.Exit(1)
			}
		}
	}()

	defer close(wait)
	logger.Info("Demo Job")
	time.Sleep(time.Second * 5)
	return nil
}
