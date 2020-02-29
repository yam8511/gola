package google

import (
	errorCode "gola/app/common/errorcode"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"io/ioutil"
	"net/http"
	"strings"
)

// CheckLogin 確認登入
func CheckLogin(sid string) (isLogin bool, apiErr errorCode.Error) {
	google := bootstrap.GetAppConf().Services.Google
	url := "http://"
	if google.Secure {
		url = "https://"
	}
	url += google.IP + google.Port + "/auth/check"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		apiErr = errorCode.GetAPIError("new_http_err", err)
		logger.Danger("CheckLogin: 建立連線請求失敗, " + err.Error())
		return
	}

	req.Header.Set("Api-Key", google.APIKey)
	req.AddCookie(&http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: sid,
	})

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		apiErr = errorCode.GetAPIError("do_request_err", err)
		logger.Danger("CheckLogin: 連線請求失敗, " + err.Error())
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		apiErr = errorCode.GetAPIError("google_api_err", err)
		logger.Danger("CheckLogin: 讀取回傳資料失敗, " + err.Error())
		return
	}

	checkText := strings.TrimSpace(string(body))
	isLogin = checkText == "Y"
	return
}
