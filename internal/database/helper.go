package database

import "fmt"

// 取連線字串
func getConnectName(driver, host, port, database, username, password string) string {
	switch driver {
	case "mysql":
		return username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8&parseTime=True"
	case "postgres":
		return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)
	}
	return ""
}
