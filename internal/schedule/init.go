package schedule

import (
	"fmt"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// 排程資訊
type crontabConf struct {
	Jobs []*CronJob `mapstructure:"job"`
}

func loadSchedule() []*CronJob {

	var jobs []*CronJob

	appRoot := bootstrap.GetAppRoot()
	appEnv := bootstrap.GetAppEnv()
	appSite := bootstrap.GetAppSite()

	vip := viper.New()
	vip.AddConfigPath(appRoot + "/config/schedule/" + appEnv)

	hasLoad := false
	fatal := func(err error) {
		logger.Error(err)
		os.Exit(1)
	}

	vip.SetConfigName(bootstrap.Default)
	if err := vip.ReadInConfig(); err == nil {
		var defaultJobsConf crontabConf
		if err := vip.Unmarshal(&defaultJobsConf); err != nil {
			filename := vip.ConfigFileUsed()
			fatal(fmt.Errorf(
				"載入 %s 錯誤： %s",
				color.HiYellowString(filename), err.Error(),
			))
		}

		jobs = append(jobs, defaultJobsConf.Jobs...)
		hasLoad = true
	}

	if appSite != bootstrap.Default {
		vip.SetConfigName(appSite)
		if err := vip.ReadInConfig(); err == nil {
			var siteJobsConf crontabConf
			if err := vip.Unmarshal(&siteJobsConf); err != nil {
				filename := vip.ConfigFileUsed()
				fatal(fmt.Errorf(
					"載入 %s 錯誤： %s",
					color.HiYellowString(filename), err.Error(),
				))
			}

			jobs = append(jobs, siteJobsConf.Jobs...)
			hasLoad = true
		}
	}

	// 假如都沒有載入檔案，噴錯誤
	if !hasLoad {
		fatal(fmt.Errorf("Schedule 沒有載入成功，請確認設定檔是否存在！"))
	}

	// 檢查背景是否有重複或是nil的
	checkJob := map[string]*CronJob{}
	for _, job := range jobs {
		if job != nil {
			_, ok := checkJob[job.Name]
			if ok {
				fatal(fmt.Errorf("Schedule 發現有重複定義的背景，請確認！ job.name只能唯一"))
			}
			checkJob[job.Name] = job
		}
	}

	jobs = []*CronJob{}
	for _, job := range checkJob {
		jobs = append(jobs, job)
	}

	return jobs
}
