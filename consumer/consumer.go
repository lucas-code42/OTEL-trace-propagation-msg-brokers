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

	var counter int
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

			fmt.Printf(`[*] event consume {"body": "%s", "header": "%v"}`, string(m.Body), m.Headers)
			fmt.Println()

			counter++
			if counter >= 5 {
				break
			}
		}
	}()
}
