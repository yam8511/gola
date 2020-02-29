package errorcode

import (
	"fmt"
	"gola/internal/logger"
	"runtime/debug"
	"strings"
)

var mappingAPIError = map[string]APIError{
	"ping":  {0, "pong", false},          // ping pong
	"panic": {500, "服務異常，請通知系統人員", true}, // 意外錯誤！

	/**
	 * Mask 相關錯誤
	 */
	"mask_api_error":      {2001, "請確認是否有連接網路，或者政府伺服器目前有問題，請稍後再試試看", false}, // 抓口罩數量API失敗
	"read_mask_csv_error": {2001, "政府伺服器給的資料有問題，等待政府資料恢復正常 (API)", false},   // 抓口罩數量API失敗

	/**
	 * Auth 相關錯誤
	 */
	"no_cookie":         {1000, "請重新登入", false},       // 取Cookie失敗
	"new_http_err":      {1001, "服務異常，請通知系統人員", true}, // 呼叫google的驗證登入，建立連線失敗
	"do_request_err":    {1002, "服務繁忙，請稍後重新登入", true}, // 呼叫google的驗證登入，連線失敗
	"google_api_err":    {1003, "服務繁忙，請稍後重新登入", true}, // 呼叫google的驗證登入，讀取回傳資料有問題
	"not_login":         {1004, "請重新登入", false},
	"permission_denied": {1005, "權限不足", false},

	/**
	 * 共通錯誤
	 */
	"param_required":        {9900, "缺少輸入參數", false},
	"param_invalid":         {9901, "輸入資料格式錯誤", false},
	"parse_err":             {9902, "資料解析錯誤", true},
	"seed_err":              {9904, "服務異常，請通知系統人員", true},
	"gorm_pool_is_timeout":  {9903, "Gorm連線池逾時", false},
	"gorm_pool_is_closed":   {9904, "Gorm連線池已經關閉", true},
	"gorm_pool_no_config":   {9904, "Gorm連線池尚未設定", true},
	"redis_pool_is_timeout": {9903, "Redis連線池逾時", false},
	"redis_pool_is_closed":  {9905, "Redis連線池已經關閉", true},
	"redis_pool_no_config":  {9905, "Redis連線池尚未設定", true},
	"undefined":             {9999, "Undefined Error", true},
}

// newAPIError 取API錯誤訊息
func newAPIError(text string, err error) Error {
	text = strings.TrimSpace(text)
	apiErroCode, ok := mappingAPIError[text]
	if !ok {
		apiErroCode = APIError{
			9999,
			fmt.Sprintf("未定義錯誤訊息 [ %s ]", text),
			true,
		}
	}

	if err != nil {
		apiErroCode.Text = fmt.Sprintf("%s (%s)", apiErroCode.Text, err.Error())
	}

	if apiErroCode.showLog {
		stack := string(debug.Stack())
		logger.Danger(fmt.Sprintf(
			"🚧  🚧  🚧  \n%s\n🎃  內部發生錯誤 [%s], %s 🎃\n🚒  🚒  🚒\n",
			stack,
			text,
			apiErroCode.Error(),
		))
	}

	return &apiErroCode
}
