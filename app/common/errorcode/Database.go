package errorcode

import (
	"gola/internal/database"
)

// CheckGormConnError 確認DB連線錯誤
func CheckGormConnError(key string, err error) (apiErr Error) {
	if err == database.ErrPoolHasNoConf {
		apiErr = newAPIError("gorm_pool_no_config", err)
		return
	}

	var yes bool
	yes = database.IsPoolTimeout(err)
	if yes {
		apiErr = newAPIError("gorm_pool_is_timeout", err)
		return
	}

	yes = database.IsPoolClosed(err)
	if yes {
		apiErr = newAPIError("gorm_pool_is_closed", err)
		return
	}

	apiErr = newAPIError(key, err)
	return
}

// CheckRedisConnError 確認快取連線錯誤
func CheckRedisConnError(key string, err error) (apiErr Error) {
	if err == database.ErrPoolHasNoConf {
		apiErr = newAPIError("redis_pool_no_config", err)
		return
	}

	var yes bool
	yes = database.IsPoolTimeout(err)
	if yes {
		apiErr = newAPIError("redis_pool_is_timeout", err)
		return
	}

	yes = database.IsPoolClosed(err)
	if yes {
		apiErr = newAPIError("redis_pool_is_closed", err)
		return
	}

	apiErr = newAPIError(key, err)
	return
}
