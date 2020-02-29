package schedule

import (
	"fmt"
	"gola/app/console"
	"sync"
	"time"

	cron "gopkg.in/robfig/cron.v2"
)

// CronJob 排程背景
type CronJob struct {
	Name          string          `mapstructure:"name"`        // 背景名稱
	Spec          string          `mapstructure:"spec"`        // 執行週期
	Cmd           string          `mapstructure:"cmd"`         // 執行工作
	IsOverlapping bool            `mapstructure:"overlapping"` // 是否可以重複
	Note          string          `mapstructure:"note"`        // 說明
	entryID       cron.EntryID    // EntryID
	running       bool            // running
	mux           *sync.RWMutex   // 讀寫鎖
	wg            *sync.WaitGroup // 等待通道
}

// Run 執行背景
func (c *CronJob) Run() {
	// logger.Info(fmt.Sprintf("開始執行 ------> %s\n", c.Name))

	c.mux.RLock()
	overlapping := c.IsOverlapping
	running := c.running
	c.mux.RUnlock()

	job := console.CronJob{
		Name: c.Name,
		Cmd:  c.Cmd,
		Spec: c.Spec,
		Note: c.Note,
	}
	// 記錄執行結束時間
	canDo := console.CanJobWork(job)

	if !canDo {
		return
	}

	// 如果不可重複，而且已經執行則跳過
	if running && !overlapping {
		// logger.Info(fmt.Sprintf("還在執行中 ------> %s\n", c.Name))
		return
	}

	// 執行背景
	c.wg.Add(1)
	c.mux.Lock()
	c.running = true
	c.mux.Unlock()
	startTime := time.Now()
	execErr := c.Exec()
	endTime := time.Now()
	c.mux.Lock()
	c.running = false
	c.mux.Unlock()
	c.wg.Done()

	// 記錄執行結束時間
	console.RecordJobStatus(job, startTime, endTime, execErr)
}

// SetEntryID 設定 EntryID
func (c *CronJob) SetEntryID(id cron.EntryID) {
	c.entryID = id
}

// Init 初始化
func (c *CronJob) Init() {
	c.mux = new(sync.RWMutex)
	c.wg = new(sync.WaitGroup)
}

// Wait 等待結束
func (c *CronJob) Wait() {
	c.wg.Wait()
}

// Exec 執行指令
func (c *CronJob) Exec() error {
	cmd := console.GetCommand(c.Cmd)
	if cmd == nil {
		return fmt.Errorf("指令尚未註冊: %s", c.Cmd)
	}
	err := cmd.Run()
	return err
}
