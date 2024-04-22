package main

import (
	"log"

	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/rabbitmq"
)

func main() {
	body := "test body"

	for i := range 10 {
		if err := rabbitmq.New().Publish(body); err != nil {
			log.Panic(err)
		}
		log.Printf("[%d] Sent %s\n", i, body)
	}
}
