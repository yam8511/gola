package handler

import (
	datastruct "gola/app/common/data_struct"
	errorCode "gola/app/common/error_code"
	"gola/app/model"
	"gola/internal/bootstrap"
	"gola/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API 範例
// @Summary 範例
// @Description 範例
// @Tags 範例
// @Success 200 {string} string "hello world"
// @Router /api/hello [get]
func API(c *gin.Context) {
	c.JSON(http.StatusOK, datastruct.API{
		Result: "hello world",
	})
}

// Seed 種子
// @Summary 種子
// @Description 在DB產生資料
// @Tags 範例
// @Success 200 {object} errorcode.APIError
// @Router /api/seed [post]
func Seed(c *gin.Context) {
	db, err := database.NewOrmConnection(true)
	if err != nil {
		bootstrap.WriteLog("ERROR", "Seed: DB連線失敗, "+err.Error())
		errorCode.GetAPIError(500)
		return
	}
	defer db.Close()

	err = model.UserSeed(db)
	if err != nil {
		bootstrap.WriteLog("ERROR", "Seed: 新增使用者資料失敗, "+err.Error())
		errorCode.GetAPIError(500)
		return
	}
}
