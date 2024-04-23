package main

import (
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/consumer"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/producer"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/script"
)

func main() {
	script.Script()
	go producer.Producer()
	go consumer.Consumer()
	select {}
}
