package bootstrap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

// 取Config檔案名稱
func GetConfigFilename() string {
	filename := strings.TrimSpace(os.Getenv("GOLA_CONFIG"))
	if filename == "" {
		filename = filepath.Join(
			GetAppRoot(),
			"config/project",
			GetAppEnv(),
			GetAppSite()+".toml",
		)
	}

	return filename
}

// GetAppConf 取專案的設定檔
func GetAppConf() Config {
	conf := defaultConf

	if vip != nil {
		_ = vip.Unmarshal(&conf)
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
