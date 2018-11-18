package google

import (
	errorCode "gola/app/common/error_code"
	"gola/internal/bootstrap"
	"io/ioutil"
	"net/http"
	"strings"
)

// CheckLogin 確認登入
func CheckLogin(sid string) (isLogin bool, apiErr errorCode.APIError) {
	google := bootstrap.GetAppConf().Servers.Google
	url := "http://"
	if google.Secure {
		url = "https://"
	}
	url += google.IP + google.Port + "/auth/check"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		apiErr = errorCode.GetAPIError(1001)
		bootstrap.WriteLog("ERROR", "CheckLogin: 建立連線請求失敗, "+err.Error())
		return
	}

	req.Header.Set("Api-Key", google.APIKey)
	req.AddCookie(&http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: sid,
	})

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		apiErr = errorCode.GetAPIError(1002)
		bootstrap.WriteLog("ERROR", "CheckLogin: 連線請求失敗, "+err.Error())
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		apiErr = errorCode.GetAPIError(1003)
		bootstrap.WriteLog("ERROR", "CheckLogin: 讀取回傳資料失敗, "+err.Error())
		return
	}

	checkText := strings.TrimSpace(string(body))
	isLogin = checkText == "Y"
	return
}
