package global

import "os"

var AppVersion string // 專案版本號
var GoVersion string  // Go版本號
var BuildTime string  // 編譯時間

func init() {
	AppVersion = os.Getenv("VERSION")
}
