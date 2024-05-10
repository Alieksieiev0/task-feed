package broker

import (
	"github.com/wagslane/go-rabbitmq"
)

func NewRabbitMQConsumer(
	conn *rabbitmq.Conn,
	queue string,
	callback func([]byte) error,
) (*RabbitMQConsumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		conn,
		queue,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsConcurrency(10),
		rabbitmq.WithConsumerOptionsQOSPrefetch(100),
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{consumer: consumer, callback: callback}, nil
}

type RabbitMQConsumer struct {
	consumer *rabbitmq.Consumer
	callback func([]byte) error
}

func (r *RabbitMQConsumer) Run(errChan chan<- error) error {
	return r.consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		errChan <- r.callback(d.Body)
		return rabbitmq.Ack
	})
}

func (r *RabbitMQConsumer) Close() {
	r.consumer.Close()
}
