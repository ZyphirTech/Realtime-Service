package main

import (
	"log"
	"net/http"

	"realtime/internal/nats"
	"realtime/internal/ws"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	hub := ws.NewHub()

	err := nats.Suscribe(hub)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		if userId == "" {
			http.Error(w, "userId is required", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		hub.Add(userId, conn)

		go func() {
			defer hub.Remove(userId)

			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	})

	log.Println("Realtime on :8081")
	http.ListenAndServe(":8081", nil)
}
