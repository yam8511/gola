package errorcode

import (
	"gola/internal/database"
)

// CheckGormConnError 確認DB連線錯誤
func CheckGormConnError(code Code, err error) (apiErr Error) {
	if err == nil {
		return
	}

	if err == database.ErrPoolHasNoConf {
		apiErr = Code_DB_No_Config.New("資料庫資訊尚未設定")
		return
	}

	var yes bool
	yes = database.IsPoolTimeout(err)
	if yes {
		apiErr = Code_DB_Timeout.New("資料庫連線池逾時")
		return
	}

	yes = database.IsPoolClosed(err)
	if yes {
		apiErr = Code_DB_Closed.New("資料庫連線池已經關閉")
		return
	}

	apiErr = code.New(err.Error())
	return
}

// CheckRedisConnError 確認快取連線錯誤
func CheckRedisConnError(code Code, err error) (apiErr Error) {
	if err == nil {
		return
	}

	if err == database.ErrPoolHasNoConf {
		apiErr = Code_DB_No_Config.New("Redis資訊尚未設定")
		return
	}

	var yes bool
	yes = database.IsPoolTimeout(err)
	if yes {
		apiErr = Code_DB_Timeout.New("Redis連線池逾時")
		return
	}

	yes = database.IsPoolClosed(err)
	if yes {
		apiErr = Code_DB_Closed.New("Redis連線池已經關閉")
		return
	}

	apiErr = code.New(err.Error())
	return
}
