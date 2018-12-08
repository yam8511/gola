package schedule

import (
	"fmt"
	"gola/app/console"
	"sync"

	cron "gopkg.in/robfig/cron.v2"
)

// CronJob 排程背景
type CronJob struct {
	// 背景名稱
	Name string `toml:"name"`
	// 執行週期
	Spec string `toml:"spec"`
	// 執行工作
	Cmd string `toml:"cmd"`
	// 是否可以重複
	IsOverlapping bool `toml:"overlapping"`
	// 說明
	Note string `toml:"note"`
	// EntryID
	entryID cron.EntryID
	// running
	running bool
	// 讀寫鎖
	mux *sync.RWMutex
	// 等待通道
	wg *sync.WaitGroup
}

// Run 執行背景
func (c *CronJob) Run() {
	// global.WriteLog("INFO", fmt.Sprintf("開始執行 ------> %s\n", c.Name))
	// TODO: 檢查背景是否可以啟動

	c.mux.RLock()
	overlapping := c.IsOverlapping
	running := c.running
	c.mux.RUnlock()

	// 如果可以重複，直接執行
	if overlapping {
		c.wg.Add(1)
		c.Exec()
		c.wg.Done()
		return
	}

	// 如果不可重複，而且已經執行則跳過
	if running {
		// global.WriteLog("INFO", fmt.Sprintf("還在執行中 ------> %s\n", c.Name))
		return
	}

	// 執行背景
	c.wg.Add(1)
	c.mux.Lock()
	c.running = true
	c.mux.Unlock()
	c.Exec()
	c.mux.Lock()
	c.running = false
	c.mux.Unlock()
	c.wg.Done()
	// TODO: 記錄執行結束時間
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
