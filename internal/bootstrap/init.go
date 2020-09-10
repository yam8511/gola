package bootstrap

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var defaultVIP *viper.Viper
var appVIP *viper.Viper

// LoadConfig 載入 config
func LoadConfig() {
	appRoot := GetAppRoot()
	appEnv := GetAppEnv()
	appSite := GetAppSite()

	replacer := strings.NewReplacer(
		".", "_",
		"SERVER.PORT", "PORT",
	)
	configDir := filepath.Join(appRoot, "config", "project", appEnv)
	conf := defaultConf

	defaultVIP = viper.New()
	defaultVIP.AutomaticEnv()
	defaultVIP.SetEnvKeyReplacer(replacer)
	defaultVIP.AddConfigPath(configDir)
	defaultVIP.SetConfigType("toml")
	defaultVIP.SetConfigName(Default)

	appVIP = viper.New()
	appVIP.AutomaticEnv()
	appVIP.SetEnvKeyReplacer(replacer)
	appVIP.AddConfigPath(configDir)
	appVIP.SetConfigType("toml")
	appVIP.SetConfigName(appSite)

	// 先讀取預設檔
	if err := defaultVIP.ReadInConfig(); err == nil {
		filename := defaultVIP.ConfigFileUsed()
		if err := defaultVIP.Unmarshal(&conf); err != nil {
			FatalLoad(filename, err)
		} else {
			log.Println(color.HiCyanString("〖 GOLA 〗[INFO] 讀取設定檔: " + filename))
		}

		defaultVIP.WatchConfig()
	}

	// 再讀取指定Site檔
	if appSite != Default {
		if err := appVIP.ReadInConfig(); err == nil {
			filename := appVIP.ConfigFileUsed()
			if err := appVIP.Unmarshal(&conf); err != nil {
				filename := appVIP.ConfigFileUsed()
				FatalLoad(filename, err)
			} else {
				log.Println(color.HiCyanString("〖 GOLA 〗[INFO] 讀取設定檔: " + filename))
			}

			appVIP.WatchConfig()
		}
	}

	appVIP.Set("app.site", appSite)
	appVIP.SetDefault("bot.token", conf.Bot.Token)
	appVIP.SetDefault("bot.chat_id", conf.Bot.ChatID)
	appVIP.SetDefault("bot.debug", conf.Bot.Debug)
}

// SetRunMode 執行模式
func SetRunMode(mode Mode) {
	if appVIP != nil {
		appVIP.Set(modeKey, mode)
	}
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
