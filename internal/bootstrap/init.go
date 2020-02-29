package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// 全域設定變數
var conf *Config

// Default 預設值
const Default = "default"

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
		if err := vip.Unmarshal(&conf); err != nil {
			filename := vip.ConfigFileUsed()
			log.Fatalf(
				"〖ERROR〗❌ 載入 %s 錯誤： %s ❌\n",
				color.HiYellowString(filename), color.HiRedString(err.Error()),
			)
		}

		hasLoad = true
	}

	// 再讀取指定Site檔
	if appSite != Default {
		vip.SetConfigName(appSite)
		if err := vip.ReadInConfig(); err == nil {
			if err := vip.Unmarshal(&conf); err != nil {
				filename := vip.ConfigFileUsed()
				log.Fatalf(
					"〖ERROR〗❌ 載入 %s 錯誤： %s ❌\n",
					color.HiYellowString(filename), color.HiRedString(err.Error()),
				)
			}
			hasLoad = true
		}
	}

	// 假如都沒有載入檔案，噴錯誤
	if !hasLoad {
		log.Fatalln(color.HiYellowString(
			"〖ERROR〗❌ Config 沒有載入成功，請確認設定檔是否存在! ❌",
		))
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

// WaitFunc 等待Func完成
func WaitFunc(fn func()) context.Context {
	ctx, finish := context.WithCancel(context.Background())
	go func() {
		fn()
		finish()
	}()
	return ctx
}
