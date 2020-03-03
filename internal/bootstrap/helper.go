package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// GetAppEnv 取環境變數
func GetAppEnv() string {
	env := strings.TrimSpace(os.Getenv("APP_ENV"))
	if env == "" {
		env = "default"
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
func GetAppConf() *Config {
	if conf != nil {
		return conf
	}
	return LoadConfig()
}
