package google

import (
	"gola/app/common/errorcode"
	"gola/internal/bootstrap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// CheckLogin 確認登入
func CheckLogin(sid string) (isLogin bool, apiErr errorcode.Error) {
	google := bootstrap.GetAppConf().Services.Google
	url := "http://"
	if google.Secure {
		url = "https://"
	}
	url += google.IP + google.Port + "/auth/check"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		apiErr = errorcode.Code_New_Request.New("CheckLogin: 建立連線請求失敗: %w", err)
		return
	}

	req.Header.Set("Api-Key", google.APIKey)
	req.AddCookie(&http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: sid,
	})

	client := http.Client{}
	if bootstrap.RunMode() == bootstrap.ServerMode {
		client.Timeout = time.Second * 30
	}
	res, err := client.Do(req)
	if err != nil {
		apiErr = errorcode.Code_Do_Request.New("CheckLogin: 連線請求失敗: %w", err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		apiErr = errorcode.Code_Google_API_Return.New("CheckLogin: 讀取回傳資料失敗: %w", err)
		return
	}

	checkText := strings.TrimSpace(string(body))
	isLogin = checkText == "Y"
	return
}
