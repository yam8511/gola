package errorcode

import (
	"fmt"
)

var mappingAPIError = map[int]APIError{
	500: {500, "服務異常，請通知系統人員"}, // 意外錯誤！

	/**
	 * Auth 相關錯誤
	 */
	1000: {1000, "請重新登入"},        // 取Cookie失敗
	1001: {1001, "服務異常，請通知系統人員"}, // 呼叫google的驗證登入，建立連線失敗
	1002: {1002, "服務繁忙，請稍後重新登入"}, // 呼叫google的驗證登入，連線失敗
	1003: {1003, "服務繁忙，請稍後重新登入"}, // 呼叫google的驗證登入，讀取回傳資料有問題
	1004: {1004, "請重新登入"},        // 尚未登入

	/**
	 * 共通錯誤
	 */
	9900: {9900, "缺少輸入參數"},          // 缺少輸入參數
	9901: {9901, "輸入資料格式錯誤"},        // 輸入參數錯誤
	9999: {9999, "Undefined Error"}, // 未定義的錯誤代碼
}

// GetAPIError 由錯誤碼取得API回傳
func GetAPIError(code int) APIError {
	if code == 0 {
		return APIError{}
	}
	api, ok := mappingAPIError[code]
	if !ok {
		return APIError{9999, fmt.Sprintf("Undefined Error (%d)", code)}
	}
	return api
}
