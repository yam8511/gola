package bootstrap

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

// GetAppEnv 取環境變數
func GetAppEnv() string {
	env := strings.TrimSpace(os.Getenv("APP_ENV"))
	if env == "" {
		env = "local"
	}
	return env
}

// GetAppSite 取賽程控客端變數
func GetAppSite() string {
	site := strings.TrimSpace(os.Getenv("APP_SITE"))
	if site == "" {
		site = "default"
	}
	return site
}

// GetAppRoot 取專案的根目錄
func GetAppRoot() string {
	var root string
	if os.Getenv("APP_ROOT") == "" {
		execRoot, err := os.Getwd()
		if err != nil {
			color.Yellow(fmt.Sprintf("[WARN] 🎃  GetAppRoot 取根目錄失敗 (%v) 🎃", err))
		}
		root = execRoot
	} else {
		root = os.Getenv("APP_ROOT")
	}

	return root
}

// GetAppConf 取專案的設定檔
func GetAppConf() Config {
	conf := defaultConf

	if defaultVIP != nil {
		_ = defaultVIP.Unmarshal(&conf)
	}

	if appVIP != nil {
		_ = appVIP.Unmarshal(&conf)
		conf.mode = Mode(appVIP.GetString(modeKey))
	}

	if conf.mode != ServerMode {
		conf.mode = CommandMode
	}
	return conf
}

// FatalLoad 載入錯誤
func FatalLoad(filename string, err error) {
	log.Fatalf(
		"〖ERROR〗❌ 載入 %s 失敗： %s ❌\n",
		color.HiYellowString(filename), color.HiRedString(err.Error()),
	)
}
