package schedule

import (
	"fmt"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"os"

	"github.com/fatih/color"
	cron "gopkg.in/robfig/cron.v2"
)

// Run 啟動排程
func Run() {
	jobs := loadSchedule()

	if len(jobs) == 0 {
		logger.Success("🎃  無定義排程，結束程序 🎃")
		return
	}

	bg := cron.New()
	for _, job := range jobs {
		job.Init()
		pid, err := bg.AddJob(job.Spec, job)
		if err != nil {
			logger.Error(fmt.Errorf(
				"%s 加入排程失敗: %s",
				color.HiYellowString(job.Name), err.Error(),
			))
		} else {
			job.SetEntryID(pid)
		}
	}

	// 開始排程
	logger.Success("🐳  排程開始啟動 🐳")
	bg.Start()

	// 等待結束訊號
	<-bootstrap.GracefulDown()
	logger.Warn("🚦  排程收到訊號囉，等待其他背景完成，準備結束排程 🚦")

	// 停止排程
	bg.Stop()

	select {
	case <-bootstrap.WaitFunc(func() {
		// 等待背景結束
		for _, job := range jobs {
			job.Wait()
		}
	}).Done():
	case <-bootstrap.WaitOnceSignal():
		logger.Danger(`🚦  收到第二次訊號，強制結束 🚦`)
		os.Exit(2)
	}

	logger.Success("🔥  排程結束囉 🔥")
}
