package supervisord

import (
	"context"
	"fmt"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"gola/internal/server"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// 執行守護進程
func DaemonExec(rootCmd *cobra.Command) {
	// rootCmd.PersistentFlags().BoolP("daemon", "d", false, "執行守護進程")
	daemon, err := rootCmd.PersistentFlags().GetBool("daemon")
	if err != nil || !daemon {
		return
	}

	restartCh := make(chan struct{})
	ExpressListen(restartCh)

	args := make([]string, 0, len(os.Args))
	for _, arg := range os.Args {
		if arg != "--daemon" && arg != "-d" {
			args = append(args, arg)
		}
	}
	bootstrap.SetupGracefulSignal()
	exit := false
	for !exit {
		ctx, done := context.WithCancel(context.Background())
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Start()
		if err != nil {
			logger.Warn(fmt.Sprintf("cmd %v start error: %s", args, err.Error()))
			done()
			continue
		}

		go func(cmd *exec.Cmd, done func()) {
			err = cmd.Wait()
			if err != nil {
				logger.Warn(fmt.Sprintf("cmd %v wait error: %s", args, err.Error()))
			}
			done()
		}(cmd, done)

		select {
		case <-ctx.Done():
		case <-restartCh:
			err = cmd.Process.Signal(syscall.SIGINT)
			if err != nil {
				logger.Warn(fmt.Sprintf("cmd %v signal error: %s", args, err.Error()))
			}
		case <-bootstrap.GracefulDown():
			exit = true
			err = cmd.Process.Signal(syscall.SIGINT)
			if err != nil {
				logger.Warn(fmt.Sprintf("cmd %v signal error: %s", args, err.Error()))
			}
		}

		select {
		case <-ctx.Done():
		case <-bootstrap.WaitOnceSignal():
			logger.Danger("⛔️  收到第二次訊號，強制結束 ⛔️")
		}

		logger.Info(fmt.Sprintf("cmd %v exit [%d] %s", args, cmd.ProcessState.ExitCode(), cmd.ProcessState.String()))
	}

	os.Exit(0)
}

func ExpressListen(restartCh chan struct{}) {
	originPORT := os.Getenv("PORT")
	os.Setenv("PORT", "4002")
	bootstrap.LoadConfig()
	os.Setenv("PORT", originPORT)
	go func() {
		server.Run(func(r *gin.Engine) {
			r.Any("/api/supervisor/restart", func(c *gin.Context) {
				select {
				case <-time.After(time.Second * 10):
					c.String(http.StatusOK, "重啟逾時，請稍候重試")
				case restartCh <- struct{}{}:
					c.String(http.StatusOK, "重啟成功")
				}
			})
		})
	}()
	time.Sleep(time.Millisecond)
}
