package utils

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// AcceptWebSocket upgrades an HTTP connection to a WebSocket connection.
func AcceptWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket:", err)
		return nil, err
	}
	return conn, nil
}
