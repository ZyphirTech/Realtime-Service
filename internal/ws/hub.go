package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mux     sync.Mutex
	Clients map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) Add(id string, conn *websocket.Conn) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.Clients[id] = conn
}

func (h *Hub) Remove(id string) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if c, ok := h.Clients[id]; ok {
		c.Close()
		delete(h.Clients, id)
	}
}

func (h *Hub) SendTo(id string, msg []byte) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if c, ok := h.Clients[id]; ok {
		err := c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			delete(h.Clients, id)
			c.Close()
		}
	}
}
