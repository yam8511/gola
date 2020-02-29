package handler

import (
	"net/http"
	"strings"

	dance "gola/criminal_dance"
	"gola/werewolf"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WerewolfGame 狼人殺
// @Summary 狼人殺
// @Tags 遊戲
// @Success 200 {string} string "狼人殺"
// @Router /api/wf/game [get]
func WerewolfGame(c *gin.Context) {
	token := strings.TrimSpace(c.Query("token"))

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	werewolf.EnterGame(conn, token)
}

// CriminalDanceGame 犯人在跳舞
// @Summary 犯人在跳舞
// @Tags 遊戲
// @Success 200 {string} string "犯人在跳舞"
// @Router /api/wf/game [get]
func CriminalDanceGame(c *gin.Context) {
	token := strings.TrimSpace(c.Query("token"))

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	dance.EnterGame(conn, token)
}
