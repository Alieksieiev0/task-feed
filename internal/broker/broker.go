package broker

import (
	"log"

	"github.com/Alieksieiev0/task-feed/internal/app"
	"github.com/wagslane/go-rabbitmq"
)

func NewRabbitMQBroker(
	storage app.Storage[func(*rabbitmq.ConnectionOptions), *rabbitmq.Conn],
	queue string,
	callback func([]byte) error,
) (*RabbitMQBroker, error) {
	conn, err := storage.Open(rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		return nil, err
	}

	consumer, err := rabbitmq.NewConsumer(
		conn,
		queue,
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsConcurrency(10),
		rabbitmq.WithConsumerOptionsQOSPrefetch(100),
	)
	if err != nil {
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQBroker{consumer: consumer, publisher: publisher, callback: callback}, nil
}

type RabbitMQBroker struct {
	publisher *rabbitmq.Publisher
	consumer  *rabbitmq.Consumer
	callback  func([]byte) error
}

func (r *RabbitMQBroker) Publish(message []byte) error {
	return r.publisher.Publish(
		message,
		[]string{"my_routing_key"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
}

func (r *RabbitMQBroker) Consume() error {
	return r.consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		err := r.callback(d.Body)
		if err != nil {
			log.Println("Consume Error: %w", err)
		}
		return rabbitmq.Ack
	})
}

func (r *RabbitMQBroker) Close() {
	r.consumer.Close()
	r.publisher.Close()
}
