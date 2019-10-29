package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Suggest 建議
// @Summary 建議
// @Description 建議
// @Tags Talk
// @Success 200 {object} errorcode.APIError
// @Router /api/suggest [post]
func Suggest(c *gin.Context) {
	token := os.Getenv("TG_TOKEN")
	chatID := os.Getenv("TG_CHAT_ID")
	if token == "" || chatID == "" {
		c.String(http.StatusOK, "failed: TELEGRAM ENV not assigned")
		return
	}

	type requestParams struct {
		Email   string `json:"email"`
		Suggest string `json:"suggest"`
		Game    string `json:"game"`
	}

	reqParams := requestParams{}
	err := c.ShouldBindJSON(&reqParams)
	if err != nil {
		c.String(http.StatusOK, "failed: "+err.Error())
		return
	}

	var gameName string
	switch reqParams.Game {
	case "wf":
		gameName = "狼人殺🐺"
	case "cd":
		gameName = "犯人在跳舞💃🕺"
	default:
		return
	}

	reqParams.Email = strings.TrimSpace(reqParams.Email)
	reqParams.Suggest = strings.TrimSpace(reqParams.Suggest)
	if reqParams.Suggest == "" {
		c.String(http.StatusOK, "failed: suggest is empty")
		return
	}

	text := fmt.Sprintf(`
		🕹%s 有使用者回饋囉！

		📮 %s

		📜
		%s
		📝
	`, gameName, reqParams.Email, reqParams.Suggest)

	query := url.Values{}
	query.Add("chat_id", chatID)
	query.Add("text", text)

	url := "https://api.telegram.org/bot" + token + "/sendMessage?" + query.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.String(http.StatusOK, "failed: "+err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusOK, "failed: "+err.Error())
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.String(http.StatusOK, "failed: "+err.Error())
		return
	}

	var resData struct {
		OK bool `json:"ok"`
	}

	err = json.Unmarshal(body, &resData)
	if err != nil {
		c.String(http.StatusOK, "failed: "+err.Error())
		return
	}

	if !resData.OK {
		c.String(http.StatusOK, "failed: "+string(body))
		return
	}

	c.String(http.StatusOK, "ok")
}
