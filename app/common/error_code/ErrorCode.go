package errorcode

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

var mappingAPIError = map[string]APIError{
	"ping":  {0, "pong"},           // ping pong
	"panic": {500, "服務異常，請通知系統人員"}, // 意外錯誤！

	/**
	 * Mask 相關錯誤
	 */
	"mask_api_error":      {2001, "請確認是否有連接網路，或者政府伺服器目前有問題，請稍後再試試看"}, // 抓口罩數量API失敗
	"read_mask_csv_error": {2001, "政府伺服器給的資料有問題，等待政府資料恢復正常 (API)"},   // 抓口罩數量API失敗

	/**
	 * Auth 相關錯誤
	 */
	"no_cookie":      {1000, "請重新登入"},        // 取Cookie失敗
	"new_http_err":   {1001, "服務異常，請通知系統人員"}, // 呼叫google的驗證登入，建立連線失敗
	"do_request_err": {1002, "服務繁忙，請稍後重新登入"}, // 呼叫google的驗證登入，連線失敗
	"google_api_err": {1003, "服務繁忙，請稍後重新登入"}, // 呼叫google的驗證登入，讀取回傳資料有問題
	"not_login":      {1004, "請重新登入"},        // 尚未登入

	/**
	 * 共通錯誤
	 */
	"param_required": {9900, "缺少輸入參數"},       // 缺少輸入參數
	"param_invalid":  {9901, "輸入資料格式錯誤"},     // 輸入參數錯誤
	"new_db_err":     {9902, "服務繁忙，請稍後嘗試"},   // 取DB連線錯誤
	"seed_err":       {9903, "服務異常，請通知系統人員"}, // Seed錯誤
	"parse_err":      {9904, "資料解析錯誤"},
	"undefined":      {9999, "Undefined Error"}, // 未定義的錯誤代碼
}

// GetAPIError 由錯誤碼取得API回傳
func GetAPIError(text string, err error) APIError {
	if err != nil {
		log.Printf("發生錯誤: [%s] %s\n", text, err.Error())
	}

	api, ok := mappingAPIError[text]
	if !ok {
		return APIError{9999, "Undefined Error (" + text + ")"}
	}
	return api
}

// GetGqlError 由錯誤碼取得Gql回傳
func GetGqlError(p graphql.ResolveParams, text string, err error) gqlerrors.FormattedError {
	if err != nil {
		log.Printf("GraphQL 發生錯誤: [%s] %s\n", text, err.Error())
	}

	// type FormattedError struct {
	// 	Message       string                    `json:"message"`
	// 	Locations     []location.SourceLocation `json:"locations"`
	// 	Path          []interface{}             `json:"path,omitempty"`
	// 	Extensions    map[string]interface{}    `json:"extensions,omitempty"`
	// 	originalError error
	// }

	apiErr := gqlerrors.NewFormattedError(text)
	if p.Info.Path != nil {
		apiErr.Path = []interface{}{p.Info.Path.Key}
	}

	return apiErr
}
