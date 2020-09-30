package handler

import (
	"gola/app/common/constant"
	"gola/app/common/errorcode"
	"gola/app/common/response"
	"gola/app/common/swagger"
	"gola/app/model"
	gogreet "gola/gorpc/greet"
	"gola/grpc/greet"
	"gola/internal/bootstrap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API 範例
// @Summary 範例
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

// API 範例
// @Summary 範例
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
