package bootstrap

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var vip *viper.Viper

// LoadConfig 載入 config
func LoadConfig() {
	filename := GetConfigFilename()

	_, err := os.Stat(filename)
	if err != nil {
		FatalLoad(filename, err)
	}

	replacer := strings.NewReplacer(
		".", "_",
		"SERVER.PORT", "PORT",
	)

	vip = viper.New()
	vip.AutomaticEnv()
	vip.SetEnvKeyReplacer(replacer)
	vip.SetConfigFile(filename)

	err = vip.ReadInConfig()
	if err != nil {
		FatalLoad(filename, err)
	}
	color.Cyan("載入配置檔: %s", vip.ConfigFileUsed())
	vip.WatchConfig()

	vip.Set("app.site", GetAppSite())
}

// SetRunMode 執行模式
func SetRunMode(m Mode) {
	mode = m
}

// RunMode 目前執行模式
func RunMode() Mode {
	return mode
}

var sig chan os.Signal
var serverClose chan struct{}

// SetupGracefulSignal 設定優雅關閉的信號
func SetupGracefulSignal() {
	sig = make(chan os.Signal, 1)
	serverClose = make(chan struct{})
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		close(serverClose)
	}()
}

// GracefulDown 優雅結束程式
func GracefulDown() <-chan struct{} {
	return serverClose
}

// WaitOnceSignal 等待一次的訊號
func WaitOnceSignal() (sig chan os.Signal) {
	sig = make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	return
}
