package swagger

import "gola/app/common/constant"

// ConfigRequest 設定資料請求
type ConfigRequest struct {
	// 使用者階層。 0:會員, 1:管理者, 2:系統管理者
	Level *constant.UserLevel `json:"level" example:"0"`
}
