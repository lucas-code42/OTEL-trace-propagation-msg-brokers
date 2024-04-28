package consumer

import (
	"context"
	"fmt"
	"log"

	opentelemetry "github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/otel"
	"github.com/lucas-code42/otel-trace-propagation-msg-brokers/pkg/rabbitmq"
	"go.opentelemetry.io/otel/propagation"
)

func Consumer() {
	msg, err := rabbitmq.New().Consume()
	if err != nil {
		log.Panic(err)
	}

	forever := make(chan bool)
	go func() {
		for m := range msg {
			rabbitHeader := m.Headers["traceID"].(string)

			otelHeader := map[string][]string{"Traceparent": {rabbitHeader}}
			p := propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			)
			ctx := p.Extract(context.Background(), propagation.HeaderCarrier(otelHeader))

			trace := opentelemetry.New(ctx, "consumer")
			_, span := trace.Start(ctx, "consumer")
			defer span.End()

			fmt.Printf("[*] msg body: %s\n", string(m.Body))
			fmt.Printf("[*] msg headers: %v\n", m.Headers)
			break
		}
	}()
	<-forever
}
