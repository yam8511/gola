package def

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// NewWebSocketConn New Socket連線
func NewWebSocketConn(conn *websocket.Conn) *WebSocketConn {
	wsConn := &WebSocketConn{
		Conn:    conn,
		Writer:  make(chan interface{}),
		closeCh: make(chan interface{}),
		mx:      &sync.RWMutex{},
	}

	go func(wsConn *WebSocketConn) {
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

	return wsConn
}

// WebSocketConn Socket連線
type WebSocketConn struct {
	Conn    *websocket.Conn
	Writer  chan interface{}
	closed  bool
	closeCh chan interface{}
	mx      *sync.RWMutex
}

// Close 關閉連線
func (conn *WebSocketConn) Close() {
	conn.mx.Lock()
	if !conn.closed {
		conn.closed = true
		close(conn.closeCh)
		conn.Conn.Close()
	}
	conn.mx.Unlock()
}

// Write 寫入連線
func (conn *WebSocketConn) Write(data interface{}) (closed bool) {
	conn.mx.RLock()
	closed = conn.closed
	conn.mx.RUnlock()
	if closed {
		return
	}

	select {
	case <-conn.closeCh:
		closed = true
	case conn.Writer <- data:
	}
	return
}
