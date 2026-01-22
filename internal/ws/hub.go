package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mux     sync.Mutex
	Clients map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Add(conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.Clients[conn] = true
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	delete(h.Clients, conn)
	conn.Close()
}

func (h *Hub) Boadcast(msg []byte) {
	h.mux.Lock()
	defer h.mux.Unlock()

	for c := range h.Clients {
		err := c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			delete(h.Clients, c)
			c.Close()
		}
	}
}
