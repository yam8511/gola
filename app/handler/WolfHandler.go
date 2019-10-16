package handler

import (
	"net/http"
	"strings"
	"time"

	datastruct "gola/app/common/data_struct"
	"gola/werewolf"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// JoinGame 加入遊戲
// @Summary 加入遊戲
// @Tags 狼人殺
// @Success 200 {object} object "角色選單"
// @Router /api/wf/game [get]
func JoinGame(c *gin.Context) {
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

	wsConn := datastruct.NewWebSocketConn(conn)

	go func(wsConn *datastruct.WebSocketConn) {
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
				if err != nil {
					return
				}
			case data := <-wsConn.Writer:
				err := conn.WriteJSON(data)
				if err != nil {
					return
				}
			}
		}
	}(wsConn)

	werewolf.EnterGame(wsConn, token)
}
