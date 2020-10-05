package database

import "fmt"

// 取連線字串
func getConnectName(driver, host string, port int, database, username, password string) string {
	switch driver {
	case "mysql":
		if port == 0 {
			port = 3306
		}
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Asia%%2FTaipei", // %%2F = / , Asia%2FTaipei = Asia/Taipei
			username, password, host, port, database,
		)
	case "postgres":
		if port == 0 {
			port = 5432
		}
		return fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=Asia/Taipei",
			host, port, database, username, password,
		)
	}
	return ""
}
