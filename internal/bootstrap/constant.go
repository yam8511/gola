package bootstrap

import "sync"

const Default = "default"     // 預設值
const StoragePath = "storage" // Storage資料夾名稱

var loadOnce = new(sync.Once)
var mode Mode = CommandMode

// 全域設定變數
var defaultConf = Config{
	App: AppConf{
		Name:  "default",
		Site:  "default",
		Debug: true,
	},
	Server: ServerConf{
		Port: 8000,
	},
	Log: LogConf{
		Mode:   "std+file",
		Prefix: "GOLA",
	},
}

// Mode 程序執行模式
type Mode string

// 程序執行模式
const (
	ServerMode  = Mode("server")
	CommandMode = Mode("command")
)
