package schedule

import (
	"errors"
	"gola/internal/bootstrap"
	"log"

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

	vip.SetConfigName(bootstrap.Default)
	if err := vip.ReadInConfig(); err == nil {
		var defaultJobsConf crontabConf
		filename := vip.ConfigFileUsed()
		if err := vip.Unmarshal(&defaultJobsConf); err != nil {
			bootstrap.FatalLoad(filename, err)
		} else {
			log.Println(color.HiCyanString("〖INFO〗讀取排程設定: " + filename))
		}

		jobs = append(jobs, defaultJobsConf.Jobs...)
		hasLoad = true
	}

	if appSite != bootstrap.Default {
		vip.SetConfigName(appSite)
		if err := vip.ReadInConfig(); err == nil {
			var siteJobsConf crontabConf
			filename := vip.ConfigFileUsed()
			if err := vip.Unmarshal(&siteJobsConf); err != nil {
				bootstrap.FatalLoad(filename, err)
			} else {
				log.Println(color.HiCyanString("〖INFO〗讀取排程設定: " + filename))
			}

			jobs = append(jobs, siteJobsConf.Jobs...)
			hasLoad = true
		}
	}

	// 假如都沒有載入檔案，噴錯誤
	if !hasLoad {
		bootstrap.FatalLoad("Schedule", errors.New("請確認『排程設定』是否存在！"))
	}

	// 檢查背景是否有重複或是nil的
	checkJob := map[string]*CronJob{}
	for _, job := range jobs {
		if job != nil {
			_, ok := checkJob[job.Name]
			if ok {
				bootstrap.FatalLoad("Schedule", errors.New("發現有重複定義的背景，請確認！ job.name只能唯一"))
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
