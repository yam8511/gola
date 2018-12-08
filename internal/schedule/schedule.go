package schedule

import (
	"gola/internal/bootstrap"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/naoina/toml"
	cron "gopkg.in/robfig/cron.v2"
)

// Run 啟動排程
func Run() {
	jobs, err := loadSchedule()
	if err != nil {
		log.Fatal("❌ 載入排程錯誤 --->" + err.Error() + "❌")
	}

	if len(jobs) == 0 {
		log.Println("🎃  無定義排程，結束程序 🎃")
		return
	}

	bg := cron.New()
	for _, job := range jobs {
		job.Init()
		pid, err := bg.AddJob(job.Spec, job)
		if err != nil {
			log.Fatalln(err)
		}
		job.SetEntryID(pid)
	}

	// 開始排程
	bootstrap.WriteLog("INFO", `
	🐳  啟動排程囉~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ 🐳
	`)
	bg.Start()

	// 等待結束訊號
	<-bootstrap.GracefulDown()
	bootstrap.WriteLog("WARNING", `
	🚦  收到訊號囉~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ 🚦
	`)

	// 停止排程
	bg.Stop()

	// 等待背景結束
	for _, job := range jobs {
		job.Wait()
	}

	bootstrap.WriteLog("INFO", `
	🔥  結束囉~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ 🔥
	`)
}

func loadSchedule() ([]*CronJob, error) {
	configDir := bootstrap.GetAppRoot() + "/config/schedule/" + bootstrap.GetAppEnv()
	dir, err := os.Open(configDir)
	if err != nil {
		return nil, err
	}

	var fileList []os.FileInfo
	fileList, err = dir.Readdir(-1)
	if err != nil {
		dir.Close()
		return nil, err
	}
	defer dir.Close()

	var jobs []*CronJob

	for i := range fileList {
		file := fileList[i]

		if strings.HasSuffix(file.Name(), ".toml") {
			tomlData, readFileErr := ioutil.ReadFile(configDir + "/" + file.Name())
			if readFileErr != nil {
				return nil, readFileErr
			}

			var jobsConf struct {
				Jobs []*CronJob `toml:"job"`
			}
			err := toml.Unmarshal(tomlData, &jobsConf)
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, jobsConf.Jobs...)
		}
	}

	return jobs, nil
}
