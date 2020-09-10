package bootstrap

// Config 型態
type Config struct {
	mode      Mode         // 程序執行模式
	App       AppConf      `mapstructure:"app"`
	Server    ServerConf   `mapstructure:"server"`
	Log       LogConf      `mapstructure:"log"`
	Databases DatabaseList `mapstructure:"database"`
	Caches    CacheList    `mapstructure:"cache"`
	Services  ServiceList  `mapstructure:"service"`
	Bot       BotConf      `mapstructure:"bot"`
}

// DatabaseList 資料庫清單
type DatabaseList struct {
	DefaultMaster *DatabaseConf `mapstructure:"default_master"`
	DefaultSlave  *DatabaseConf `mapstructure:"default_slave"`
}

// CacheList 快取資料庫清單
type CacheList struct {
	DefaultMaster *CacheConf `mapstructure:"default_master"`
	DefaultSlave  *CacheConf `mapstructure:"default_slave"`
}

// ServiceList 服務清單
type ServiceList struct {
	Google *ServiceConf `mapstructure:"google"`
}

// AppConf 專案資訊
type AppConf struct {
	Name  string `mapstructure:"name"`  // 專案名稱
	Env   string `mapstructure:"env"`   // 專案環境
	Site  string `mapstructure:"site"`  // 專案站別
	Debug bool   `mapstructure:"debug"` // 開啟Debug模式
	Salt  string `mapstructure:"salt"`  // 專案雜湊碼
}

// ServerConf 伺服器資訊
type ServerConf struct {
	IP      string `mapstructure:"ip"`       // 伺服器的IP
	Host    string `mapstructure:"host"`     // 伺服器的Host
	Port    int    `mapstructure:"port"`     // 伺服器的Port
	Secure  bool   `mapstructure:"secure"`   // 是否要安全憑證
	MaxConn int    `mapstructure:"max_conn"` // 最大連線數量
}

// LogConf 紀錄Log資訊
type LogConf struct {
	Mode   string `mapstructure:"mode"`   // Log紀錄模式： std, file, std+file
	Prefix string `mapstructure:"prefix"` // Log前綴
}

// DatabaseConf 資料庫設定
type DatabaseConf struct {
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	MaxConn  int    `mapstructure:"max_conn"`
}

// CacheConf 快取資料庫設定
type CacheConf struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	MaxConn  int    `mapstructure:"max_conn"`
}

// ServiceConf 服務資訊
type ServiceConf struct {
	IP     string `mapstructure:"ip"`
	Port   string `mapstructure:"port"`
	Host   string `mapstructure:"host"`
	Secure bool   `mapstructure:"secure"`
	APIKey string `mapstructure:"api_key"`
}

// BotConf 機器人設定
type BotConf struct {
	Token  string `mapstructure:"token"`
	ChatID int64  `mapstructure:"chat_id"`
	Debug  bool   `mapstructure:"debug"`
}
