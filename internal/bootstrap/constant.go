package bootstrap

// 全域設定變數
var conf *Config

// Default 預設值
const Default = "default"

// StoragePath Storage資料夾名稱
const StoragePath = "storage"

// Mode 程序執行模式
type Mode string

// 程序執行模式
const (
	ServerMode  = Mode("server")
	CommandMode = Mode("command")
)
