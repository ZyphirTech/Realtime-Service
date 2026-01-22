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

	err := nats.Suscribe(hub.Boadcast)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		hub.Add(conn)

		go func() {
			defer hub.Remove(conn)

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
