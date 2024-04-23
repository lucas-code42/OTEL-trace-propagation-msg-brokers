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
			v, b := m.Headers["traceID"]
			fmt.Println(v, b)

			p := propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			)
			x := propagation.MapCarrier{}
			ctx := p.Extract(context.Background(), x)

			fmt.Println(x.Get("traceID"))

			trace := opentelemetry.New(ctx, "consumer")
			_, span := trace.Start(ctx, "consumer")
			defer span.End()

			fmt.Printf("[*] msg body: %s\n", string(m.Body))
			fmt.Printf("[*] msg headers: %v\n", m.Headers)
			break
		}
	}()

	fmt.Println("aqui fora")
	<-forever
}
