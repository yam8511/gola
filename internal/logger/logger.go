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
func Info(text string, args ...interface{}) {
	writeLog(levelInfo, text, args...)
}

// Success 成功資訊
func Success(text string, v ...interface{}) {
	writeLog(levelSuccess, text, v...)
}

// Warn 警告資訊
func Warn(text string, v ...interface{}) {
	writeLog(levelWarn, text, v...)
}

// Danger 危險資訊
func Danger(text string, v ...interface{}) {
	writeLog(levelDanger, text, v...)
}

// Error 錯誤資訊
func Error(err error) {
	if err != nil {
		writeLog(levelError, err.Error())
	}
}

func writeLog(t level, text string, v ...interface{}) {
	conf := bootstrap.GetAppConf()
	var colorFmt *color.Color
	switch t {
	case levelInfo:
		text = "[INFO]   " + text
		colorFmt = color.New(color.FgHiBlue)
	case levelSuccess:
		text = "[OK]     " + text
		colorFmt = color.New(color.FgHiGreen)
	case levelWarn:
		text = "[WARN]   " + text
		colorFmt = color.New(color.FgHiYellow)
	case levelDanger:
		text = "[DANGER] " + text
		colorFmt = color.New(color.FgHiRed)
	case levelError:
		text = "[ERROR]  " + text
		colorFmt = color.New(color.BgHiRed, color.FgHiWhite)
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

	if conf.App.Site == "" {
		prefix = "〖 " + prefix + " 〗"
	} else {
		prefix = "〖 " + prefix + " | " + conf.App.Site + " 〗"
	}

	if conf.Log.Mode == "file" || conf.Log.Mode == "std+file" {
		w, err := createLogFile()
		logger := log.New(w, prefix, log.LstdFlags|log.Lmsgprefix)
		if err == nil {
			colorFmt.DisableColor()
			logger.Println(colorFmt.Sprintf(text, v...))
			w.Close()

			if conf.Log.Mode == "file" {
				return
			}
		} else {
			color.Yellow("[WARN]   因為打開Log檔失敗，所以直接顯示在stdout")
		}
	}

	prefix = time.Now().Format("2006-01-02 15:04:05 ") + prefix
	colorFmt.EnableColor()
	colorFmt.Printf(prefix+text+"\n", v...)
}
