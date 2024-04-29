package main

import (
	"fmt"

	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/consumer"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/producer"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/script"
)

func main() {
	fmt.Println("acess http://localhost:16686/ to check traces")
	
	script.Script()
	
	go producer.Producer()
	go consumer.Consumer()

	select {}
}
