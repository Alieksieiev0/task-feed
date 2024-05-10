package broker

import (
	"github.com/wagslane/go-rabbitmq"
)

func NewRabbitMQPublisher(conn *rabbitmq.Conn) (*RabbitMQPublisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQPublisher{publiser: publisher}, nil
}

type RabbitMQPublisher struct {
	publiser *rabbitmq.Publisher
}

func (r *RabbitMQPublisher) Run(message []byte) error {
	return r.publiser.Publish(
		message,
		[]string{"my_routing_key"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
}

func (r *RabbitMQPublisher) Close() {
	r.publiser.Close()
}
