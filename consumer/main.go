package main

import (
	"fmt"
	"log"

	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/rabbitmq"
)

func main() {
	msg, err := rabbitmq.New().Consume()
	if err != nil {
		log.Panic(err)
	}

	forever := make(chan bool)

	go func() {
		for m := range msg {
			fmt.Printf("[*] msg body: %s\n", string(m.Body))
			fmt.Printf("[*] msg headers: %s\n", string(m.Body))
		}
	}()

	<-forever
}
