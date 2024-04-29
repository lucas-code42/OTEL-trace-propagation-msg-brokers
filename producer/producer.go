package producer

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	opentelemetry "github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/otel"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/rabbitmq"
	"go.opentelemetry.io/otel/propagation"
)

func Producer() {
	ctx := context.Background()

	tracer := opentelemetry.New(ctx, "producer")
	ctx, span := tracer.Start(ctx, "producer")
	defer span.End()

	otelHeader := map[string][]string{}
	p := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	p.Inject(ctx, propagation.HeaderCarrier(otelHeader))

	for i := 0; i < 5; i++ {
		body := uuid.NewString()
		if err := rabbitmq.New().Publish(ctx, body, otelHeader); err != nil {
			log.Panic(err)
		}
		fmt.Printf("[*] event sent - %s\n", body)
	}
}
