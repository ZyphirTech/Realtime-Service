package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

func Suscribe(hub func(msg []byte)) error {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return err
	}

	_, err = nc.Subscribe("events.test", func(msg *nats.Msg) {
		log.Println("evento recibido")
		hub(msg.Data)
	})

	return err
}
