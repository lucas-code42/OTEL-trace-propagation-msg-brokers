package rabbitmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/propagation"
)

type rabbitmq struct {
	conn *amqp.Connection
}

func New() *rabbitmq {
	return &rabbitmq{
		conn: connect(),
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func connect() *amqp.Connection {
	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	return conn
}

func (r *rabbitmq) DeclareQueue(queueName string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *rabbitmq) Publish(ctx context.Context, body string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}

	// traceId := map[string]string{"traceID": "1234567890"}

	otelHeader := map[string][]string{}
	p := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	p.Inject(ctx, propagation.HeaderCarrier(otelHeader))

	header := make(amqp.Table)
	header["traceID"] = otelHeader["Traceparent"][0]

	err = ch.PublishWithContext(
		ctx,     // context
		"",      // exchange
		"hello", // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers:     header,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *rabbitmq) Consume() (<-chan amqp.Delivery, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	msg, err := ch.ConsumeWithContext(
		context.Background(),
		"hello", // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
