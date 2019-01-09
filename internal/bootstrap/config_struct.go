package bootstrap

// App 應用環境
type App struct {
	Env      string `toml:"env"`
	IP       string `toml:"ip"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Secure   bool   `toml:"secure"`
	Debug    bool   `toml:"debug"`
	AutoPort bool   `toml:"auto_port"`
}

// DBConf 資料庫設定
type DBConf struct {
	DB       string `toml:"db"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	MaxConn  int    `toml:"max_conn"`
}

// BotConf 機器人設定
type BotConf struct {
	Token  string `toml:"token"`
	ChatID int64  `toml:"chat_id"`
	Debug  bool   `toml:"debug"`
}

// Config 型態
type Config struct {
	App      App        `toml:"app"`
	Bot      BotConf    `toml:"bot"`
	DBMaster *DBConf    `toml:"database_master"`
	DBSlave  *DBConf    `toml:"database_slave"`
	Servers  ServerList `toml:"server"`
}

// ServerList 服務清單
type ServerList struct {
	Google Server `toml:"google"`
}

// Server 服務資訊
type Server struct {
	IP     string `toml:"ip"`
	Port   string `toml:"port"`
	Host   string `toml:"host"`
	Secure bool   `toml:"secure"`
	APIKey string `toml:"api_key"`
}
