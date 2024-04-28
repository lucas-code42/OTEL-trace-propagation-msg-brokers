package producer

import (
	"context"
	"log"

	opentelemetry "github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/otel"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/rabbitmq"
	"go.opentelemetry.io/otel/propagation"
)

func Producer() {
	body := "foobar"
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

	if err := rabbitmq.New().Publish(ctx, body, otelHeader); err != nil {
		log.Panic(err)
	}

	log.Printf("[*] Sent %s\n", body)
}
