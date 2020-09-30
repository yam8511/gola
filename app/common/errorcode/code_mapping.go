package errorcode

const (
	Code_Ping  = Code(0)   // ping pong
	Code_Panic = Code(500) // 意外錯誤！

	//####################################
	//             Mask 相關
	//####################################
	Code_Mask_API = Code(2001)

	//####################################
	//             Auth 相關
	//####################################

	Code_No_Cookie         = Code(1000) // 取Cookie失敗
	Code_New_Request       = Code(1001) // 建立新請求失敗
	Code_Do_Request        = Code(1002) // 進行連線失敗
	Code_Google_API_Return = Code(1003) // Google API 回傳有問題
	Code_No_Login          = Code(1004) // 未登入
	Code_Permission_Denied = Code(1005) // 權限不足

	//####################################
	//             共通 相關
	//####################################
	Code_Param_Required  = Code(9900) // 缺少輸入參數
	Code_param_Invalid   = Code(9901) // 輸入資料格式錯誤
	Code_Data_Parse      = Code(9902) // 資料解析錯誤
	Code_Seed            = Code(9904) // 自動產生資料失敗
	Code_DB_Timeout      = Code(9903) // Gorm連線池逾時
	Code_DB_Closed       = Code(9904) // Gorm連線池已經關閉
	Code_DB_No_Config    = Code(9904) // Gorm連線池尚未設定
	Code_Redos_timeout   = Code(9903) // Redis連線池逾時
	Code_Redis_Closed    = Code(9905) // Redis連線池已經關閉
	Code_Redis_No_Config = Code(9905) // Redis連線池尚未設定
	Code_Undefined       = Code(9999) // Undefined Error
)
