package websockets

import (
	"log"
	"medovukha/services/docker/containers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// !Dev Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsHub struct {
	conns      map[*websocket.Conn]bool
	connsCount int64
	broadcast  chan any
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func (h *wsHub) Run(manager *containers.DockerStreamManager) {
	for {
		select {
		case conn := <-h.register:
			{
				h.conns[conn] = true
				manager.OnClientConnected()
			}
		case conn := <-h.unregister:
			{
				if _, ok := h.conns[conn]; ok {
					delete(h.conns, conn)
					conn.Close()
				}
				manager.OnClientDisconnected()
			}
		case msg := <-h.broadcast:
			{
				for conn := range h.conns {
					if err := conn.WriteJSON(msg); err != nil {
						conn.Close()
						delete(h.conns, conn)
					}
				}
			}
		}
	}
}

func NewHub() *wsHub {
	return &wsHub{
		conns:      make(map[*websocket.Conn]bool),
		broadcast:  make(chan any),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func WsHandler(h *wsHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WsHandler error: %s\n", err.Error())
			return
		}

		h.register <- conn

		defer func() {
			conn.Close()
			h.unregister <- conn
		}()
	}
}
