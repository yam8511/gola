package handler

import (
	"gola/app/common/constant"
	"gola/app/common/errorcode"
	"gola/app/common/response"
	"gola/app/common/swagger"
	"gola/app/model"
	"gola/global"
	gogreet "gola/gorpc/greet"
	"gola/grpc/discover"
	"gola/grpc/greet"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// @Summary grpc greet
// @Description 範例
// @Tags 範例
// @Success 200 {string} string "hello world"
// @Router /api/grpc/hello [get]
func API(c *gin.Context) {
	res, err := greet.Client(c.DefaultQuery("name", "world"))
	if err != nil {
		response.Failed(c, errorcode.Code_Google_API_Return.New(err.Error()))
	} else {
		response.Success(c, res.GetMessage())
	}
}

// @Summary gorpc greet
// @Description 範例
// @Tags 範例
// @Success 200 {string} string "hello world"
// @Router /api/gorpc/hello [get]
func API2(c *gin.Context) {
	res, err := gogreet.Client(c.DefaultQuery("name", "world"))
	if err != nil {
		response.Failed(c, errorcode.Code_Google_API_Return.New(err.Error()))
	} else {
		response.Success(c, res.Message)
	}
}

// @Summary http greet
// @Description 範例
// @Tags 範例
// @Success 200 {string} string "hello world"
// @Router /api/http/hello [get]
func API3(c *gin.Context) {
	link := &url.URL{
		Scheme: "http",
		Path:   "/api/greet",
	}
	link.Host = discover.Discover("greet", "http")
	res, err := http.DefaultClient.Get(link.String())
	if err != nil {
		response.Failed(c, errorcode.Code_Google_API_Return.New(err.Error()))
		return
	}
	_ = res.Body.Close()

	if res.StatusCode != http.StatusOK {
		response.Failed(c, errorcode.Code_Google_API_Return.New(res.Status))
		return
	}

	query := link.Query()
	query.Set("again", "Y")
	link.RawQuery = query.Encode()

	res, err = http.DefaultClient.Get(link.String())
	if err != nil {
		response.Failed(c, errorcode.Code_Google_API_Return.New(err.Error()))
		return
	}

	if res.StatusCode != http.StatusOK {
		response.Failed(c, errorcode.Code_Google_API_Return.New(res.Status))
		return
	}

	defer res.Body.Close()

	msg, err := ioutil.ReadAll(res.Body)
	if err != nil {
		response.Failed(c, errorcode.Code_Google_API_Return.New(err.Error()))
		return
	}

	response.Success(c, string(msg))
}

var mx = &sync.RWMutex{}
var count int

// @Summary http greet
// @Description 範例
// @Tags 範例
// @Success 200 {string} string "hello world"
// @Router /api/greet [get]
func Greet(c *gin.Context) {
	name := c.DefaultQuery("name", "world")
	again := c.Query("again") == "Y"

	mx.Lock()
	count++
	num := strconv.Itoa(count)
	mx.Unlock()

	if again {
		msg := "#" + num + ": Get Hello Again " + global.AppVersion + " from " + name
		logger.Success(msg)
		// time.Sleep(time.Second * 2)
		c.String(http.StatusOK, msg)
		return
	}

	msg := "#" + num + ": Get Hello " + global.AppVersion + " from " + name
	logger.Success(msg)
	c.String(http.StatusOK, msg)
}

// Config 範例
// @Summary 範例
// @Description 範例
// @Tags 範例
// @Accept json
// @Produce json
// @Param  body  body  swagger.ConfigRequest   true    "使用者等級"
// @Failure 403 {object} response.API "回傳權限不足"
// @Success 200 {object} bootstrap.Config "回傳設定資料"
// @Router /api/config [post]
func Config(c *gin.Context) {
	req := swagger.ConfigRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		apiErr := errorcode.Code_param_Invalid.New("")
		response.Failed(c, apiErr)
		return
	}

	if req.Level == nil {
		apiErr := errorcode.Code_Param_Required.New("")
		response.Failed(c, apiErr)
		return
	}

	if !req.Level.IsValid() {
		apiErr := errorcode.Code_param_Invalid.New("")
		response.Failed(c, apiErr)
		return
	}

	if *req.Level == constant.SuperUserLevel {
		response.Success(c, bootstrap.GetAppConf())
		return
	}

	apiErr := errorcode.Code_Permission_Denied.New("")
	response.Failed(c, apiErr, response.WithStatusCode(http.StatusForbidden))
}

// Seed 種子
// @Summary 種子
// @Description 在DB產生資料
// @Tags 範例
// @Success 200 {object} response.API "回傳執行狀態"
// @Router /api/seed [post]
func Seed(c *gin.Context) {
	err := model.UserSeed()
	if err != nil {
		response.Failed(c, errorcode.Code_Seed.New(err.Error()))
		return
	}
	response.Success(c, nil)
}
