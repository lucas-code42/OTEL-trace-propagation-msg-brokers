package producer

import (
	"context"
	"log"

	opentelemetry "github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/otel"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/rabbitmq"
	"go.opentelemetry.io/otel/propagation"
)

func Producer() {
	body := "test body"

	ctx := context.Background()

	tracer := opentelemetry.New(ctx, "producer")
	ctx, span := tracer.Start(ctx, "producer")
	defer span.End()

	if err := rabbitmq.New().Publish(body); err != nil {
		log.Panic(err)
	}
	p := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	x := propagation.MapCarrier{}
	x.Set("traceID", "123456789")

	p.Inject(ctx, propagation.MapCarrier(x))

	log.Printf("[*] Sent %s\n", body)
}
