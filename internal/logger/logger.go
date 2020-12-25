package logger

import (
	"github.com/sirupsen/logrus"
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
	switch t {
	case levelInfo:
		logrus.WithField("lvl", levelInfo).Infof(text, v...)
	case levelSuccess:
		logrus.WithField("lvl", levelSuccess).Infof(text, v...)
	case levelWarn:
		logrus.WithField("lvl", levelWarn).Warnf(text, v...)
	case levelDanger:
		logrus.WithField("lvl", levelDanger).Errorf(text, v...)
	case levelError:
		logrus.WithField("lvl", levelError).Errorf(text, v...)
	}
}
