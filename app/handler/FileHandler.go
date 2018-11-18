package handler

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// HomeFile 首頁
func HomeFile() gin.HandlerFunc {
	return static.Serve("/", static.LocalFile("./public", true))
}
