package bootstrap

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/naoina/toml"
)

// Conf 全域設定變數
var Conf *Config

// LoadConfig 載入 config
func LoadConfig() *Config {
	configFile := GetAppRoot() + "/config/project/" + GetAppEnv() + ".toml"
	tomlData, readFileErr := ioutil.ReadFile(configFile)
	if readFileErr != nil {
		msg := fmt.Sprintf("❌ 讀取Config錯誤： %v ❌", readFileErr)
		WriteLog("ERROR", msg)
		os.Exit(1)
	}

	err := toml.Unmarshal(tomlData, &Conf)
	if err != nil {
		msg := fmt.Sprintf("❌ 載入Config錯誤： %v ❌", err)
		WriteLog("ERROR", msg)
		os.Exit(1)
	}
	return Conf
}

var sig chan os.Signal
var serverClose chan os.Signal

// SetupGracefulSignal 設定優雅關閉的信號
func SetupGracefulSignal() {
	sig = make(chan os.Signal, 1)
	serverClose = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
	go func() {
		s := <-sig
		for {
			select {
			case serverClose <- s:
			}
		}
	}()
}

// GracefulDown 優雅結束程式
func GracefulDown() <-chan os.Signal {
	return serverClose
}
