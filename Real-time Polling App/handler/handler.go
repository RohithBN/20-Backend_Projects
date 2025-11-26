package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true 
    },
}

func PollHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte("Received message ->"+string(msg))); err != nil {
			break
		}
	}
}
