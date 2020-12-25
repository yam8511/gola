package bootstrap

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	resetLogger()
}

func resetLogger() {
	conf := GetAppConf()

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
	logrus.StandardLogger().Level = logrus.TraceLevel

	if strings.Contains(conf.Log.Mode, "file") {
		logrus.AddHook(&GoLaLogHook{})
	}

	if conf.Log.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}

}

type GoLaLogHook struct{}

func (*GoLaLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func createLogFile(conf Config) (io.WriteCloser, error) {
	name := conf.App.Name
	if name == "" {
		name = "default"
	}
	name += ".log"
	now := time.Now()
	filename := GetAppRoot() + now.Format("/storage/logs/2006-01-02/15/") + name

CREATE:
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		if os.IsNotExist(err) {
			//建立資料夾
			err := os.MkdirAll(filepath.Dir(filename), 0777)
			if err != nil {
				log.Fatal("ERROR", "CreateFile: 建立相關資料夾錯誤, "+err.Error())
			}
			goto CREATE
		}
		return nil, err
	}
	return f, nil
}

func (*GoLaLogHook) Fire(lg *logrus.Entry) error {
	conf := GetAppConf()

	f, err := createLogFile(conf)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := logrus.NewEntry(logrus.New())
	entry.Logger.SetOutput(f)
	if conf.Log.Format == "json" {
		entry.Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		entry.Logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}

	entry = entry.WithFields(logrus.Fields{
		"app_name": conf.App.Name,
		"app_site": conf.App.Site,
		"app_env":  conf.App.Env,
	})

	if lg.Data != nil {
		entry = entry.WithFields(lg.Data)
	}

	switch lg.Level {
	case logrus.FatalLevel:
		entry.WithField("fatal", true).Error(lg.Message)
	case logrus.PanicLevel:
		entry.WithField("panic", true).Error(lg.Message)
	default:
		entry.Logf(lg.Level, lg.Message)
	}
	return nil
}
