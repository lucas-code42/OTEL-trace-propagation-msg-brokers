package producer

import (
	"context"
	"log"

	opentelemetry "github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/otel"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/rabbitmq"
)

func Producer() {
	body := "test body"
	ctx := context.Background()

	tracer := opentelemetry.New(ctx, "producer")
	ctx, span := tracer.Start(ctx, "producer")
	defer span.End()

	if err := rabbitmq.New().Publish(ctx, body); err != nil {
		log.Panic(err)
	}

	log.Printf("[*] Sent %s\n", body)
}
