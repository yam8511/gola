package bootstrap

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// LoadConfig 載入 config
func LoadConfig() *Config {
	appRoot := GetAppRoot()
	appEnv := GetAppEnv()
	appSite := GetAppSite()
	vip := viper.New()
	vip.AutomaticEnv()
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vip.AddConfigPath(appRoot + "/config/project/" + appEnv)

	hasLoad := false

	// 先讀取預設檔
	vip.SetConfigName(Default)
	if err := vip.ReadInConfig(); err == nil {
		filename := vip.ConfigFileUsed()
		if err := vip.Unmarshal(&conf); err != nil {
			FatalLoad(filename, err)
		} else {
			log.Println(color.HiCyanString("〖INFO〗讀取設定檔: " + filename))
		}

		hasLoad = true
	}

	// 再讀取指定Site檔
	if appSite != Default {
		vip.SetConfigName(appSite)
		if err := vip.ReadInConfig(); err == nil {
			filename := vip.ConfigFileUsed()
			if err := vip.Unmarshal(&conf); err != nil {
				filename := vip.ConfigFileUsed()
				FatalLoad(filename, err)
			} else {
				log.Println(color.HiCyanString("〖INFO〗讀取設定檔: " + filename))
			}
			hasLoad = true
		}
	}

	// 假如都沒有載入檔案，噴錯誤
	if !hasLoad {
		FatalLoad("Config", errors.New("請確認『設定檔』是否存在！"))
	}

	vip.SetDefault("bot.token", conf.Bot.Token)
	vip.SetDefault("bot.chat_id", conf.Bot.ChatID)
	vip.SetDefault("bot.debug", conf.Bot.Debug)
	conf.Bot = BotConf{
		Token:  vip.GetString("bot.token"),
		ChatID: vip.GetInt64("bot.chat_id"),
		Debug:  vip.GetBool("bot.debug"),
	}

	return conf
}

// SetRunMode 執行模式
func SetRunMode(mode Mode) {
	GetAppConf().mode = mode
}

// RunMode 目前執行模式
func RunMode() Mode {
	return GetAppConf().mode
}

var sig chan os.Signal
var serverClose chan struct{}

// SetupGracefulSignal 設定優雅關閉的信號
func SetupGracefulSignal() {
	sig = make(chan os.Signal, 1)
	serverClose = make(chan struct{})
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
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
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	return
}
