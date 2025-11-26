package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Connection struct {
	uid  string
	Conn *websocket.Conn
	mu   sync.Mutex
}

var (
	socketRoom   = make(map[string][]*Connection)
	socketRoomMu sync.RWMutex
	nextID       uint64
)

func addToRoom(room string, c *Connection) {
	socketRoomMu.Lock()
	defer socketRoomMu.Unlock()
	socketRoom[room] = append(socketRoom[room], c)
}

func removeFromRoom(room string, c *Connection) {
	socketRoomMu.Lock()
	defer socketRoomMu.Unlock()

	conns := socketRoom[room]
	for i, conn := range conns {
		if conn == c {
			socketRoom[room] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

func getRoomSnapshot(room string) []*Connection {
	socketRoomMu.RLock()
	defer socketRoomMu.RUnlock()

	conns := socketRoom[room]
	out := make([]*Connection, len(conns))
	copy(out, conns)
	return out
}

func PollHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	uid := fmt.Sprintf("%d", atomic.AddUint64(&nextID, 1))
	c := &Connection{uid: uid, Conn: conn}

	log.Println("New connection with uid:", uid)

	addToRoom("default", c)
	defer removeFromRoom("default", c)

	// Send welcome message
	c.mu.Lock()
	err = c.Conn.WriteMessage(websocket.TextMessage, []byte("Welcome to the Room"))
	c.mu.Unlock()
	if err != nil {
		log.Println("Initial write error:", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return 
		}
		log.Printf("Received from %s: %s\n", uid, message)

		// Broadcast to all other connections in the room
		conns := getRoomSnapshot("default")

		for _, other := range conns {
			if other.uid == uid {
				continue
			}
			other.mu.Lock()
			err = other.Conn.WriteMessage(websocket.TextMessage, message)
			other.mu.Unlock()
			if err != nil {
				log.Println("Write error to", other.uid, ":", err)
			}
		}
	}
}
