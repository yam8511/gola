package logger

import (
	"gola/internal/bootstrap"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

type level string

const (
	levelInfo    = level("info")
	levelSuccess = level("success")
	levelWarn    = level("warn")
	levelDanger  = level("danger")
	levelError   = level("error")
)

// Info 一般資訊
func Info(text string) {
	writeLog(levelInfo, text)
}

// Success 成功資訊
func Success(text string) {
	writeLog(levelSuccess, text)
}

// Warn 警告資訊
func Warn(text string) {
	writeLog(levelWarn, text)
}

// Danger 危險資訊
func Danger(text string) {
	writeLog(levelDanger, text)
}

// Error 錯誤資訊
func Error(err error) {
	if err != nil {
		writeLog(levelError, err.Error())
	}
}

func writeLog(t level, text string) {
	conf := bootstrap.GetAppConf()
	var coloarText string
	switch t {
	case levelInfo:
		text = "[INFO]   " + text
		coloarText = color.HiBlueString(text)
	case levelSuccess:
		text = "[OK]     " + text
		coloarText = color.HiGreenString(text)
	case levelWarn:
		text = "[WARN]   " + text
		coloarText = color.HiYellowString(text)
	case levelDanger:
		text = "[DANGER] " + text
		coloarText = color.HiRedString(text)
	case levelError:
		text = "[ERROR]  " + text
		coloarText = color.New(color.BgHiRed, color.FgHiWhite).Sprint(text)
	}

	createLogFile := func() (io.WriteCloser, error) {
		name := conf.App.Name
		if name == "" {
			name = "default"
		}
		name += ".log"
		now := time.Now()
		filename := bootstrap.GetAppRoot() + now.Format("/storage/logs/2006-01-02/15/") + name

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

	prefix := conf.Log.Prefix
	if prefix == "" {
		prefix = "GOLA"
	}
	prefix = "〖 " + prefix + " 〗"
	if conf.Log.Mode == "file" || conf.Log.Mode == "std+file" {
		w, err := createLogFile()
		logger := log.New(w, prefix, log.LstdFlags|log.Lmsgprefix)
		if err == nil {
			logger.Println(text)
			w.Close()

			if conf.Log.Mode == "file" {
				return
			}
		} else {
			color.Yellow("[WARN]   因為打開Log檔失敗，所以直接顯示在stdout")
		}
	}

	log.Printf("%s%s\n", prefix, coloarText)
}
