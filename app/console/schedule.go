package console

import (
	"time"
)

// CanJobWork 可以用來確認背景是否可以執行的Func
func CanJobWork(job CronJob) bool {
	// 可以寫自己邏輯...
	// logger.Info(job.Name + "可以執行嗎??")
	return true
}

// RecordJobStatus 背景執行完畢，會呼叫這個Func來紀錄執行狀態
func RecordJobStatus(
	job CronJob,
	startTime, endTime time.Time,
	err error,
) {
	// 可以寫自己邏輯...
	// logger.Warn(job.Name + "執行完成" + startTime.String())
}

// CronJob 背景
type CronJob struct {
	// 背景名稱
	Name string `mapstructure:"name"`
	// 執行週期
	Spec string `mapstructure:"spec"`
	// 執行工作
	Cmd string `mapstructure:"cmd"`
	// 是否可以重複
	IsOverlapping bool `mapstructure:"overlapping"`
	// 說明
	Note string `mapstructure:"note"`
}
