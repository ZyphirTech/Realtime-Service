package nats

import (
	"encoding/json"
	"log"

	"realtime/internal/ws"

	"github.com/nats-io/nats.go"
)

func Suscribe(hub *ws.Hub) error {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return err
	}

	_, err = nc.Subscribe("chat.message", func(m *nats.Msg) {
		var msg map[string]string
		json.Unmarshal(m.Data, &msg)
		log.Println(msg)
		hub.SendTo(msg["ToId"], m.Data)
	})

	return err
}
