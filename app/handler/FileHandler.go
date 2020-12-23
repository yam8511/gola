package handler

import (
	"gola/internal/logger"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// HomeFile 首頁
func HomeFile() gin.HandlerFunc {
	fs := static.Serve("/", static.LocalFile("./public", true))
	return func(c *gin.Context) {
		logger.Success("---> %s", c.Request.RequestURI)
		if c.Request.RequestURI == "/wf/" || c.Request.RequestURI == "/" {
			if pusher := c.Writer.Pusher(); pusher != nil {
				push := func(name string) {
					err := pusher.Push(name, nil)
					if err != nil {
						logger.Warn("HTTP2 Push %s 失敗: %s", name, err)
					} else {
						logger.Success("HTTP2 Push %s 成功", name)
					}
				}
				push("/images/jkopay.jpg")
				push("/images/richart.jpg")
				push(c.Request.RequestURI + "vue-qrcode.js")
			}
		}
		fs(c)
	}
}
