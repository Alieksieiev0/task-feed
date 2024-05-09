package broker

import (
	"bytes"
	"io"

	"github.com/wagslane/go-rabbitmq"
	"golang.org/x/net/context"
)

func NewRabbitMQ(addr string, queue string) (*RabbitMQ, error) {
	conn, err := rabbitmq.NewConn(addr, rabbitmq.WithConnectionOptionsLogging)
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

	ctx, cancel := context.WithCancel(context.Background())
	return &RabbitMQ{publiser: publisher, consumer: consumer, ctx: ctx, cancel: cancel}, nil
}

type RabbitMQ struct {
	publiser *rabbitmq.Publisher
	consumer *rabbitmq.Consumer
	ctx      context.Context
	cancel   context.CancelFunc
}

func (r *RabbitMQ) Run(consumeFunc func(io.Reader) error) (chan<- string, <-chan error) {
	msg := make(chan string)
	err := make(chan error)

	go r.publish(msg, err)
	go r.consume(nil, err)

	return msg, err
}

func (r *RabbitMQ) publish(in <-chan string, out chan<- error) {
	for {
		select {
		case <-r.ctx.Done():
			r.publiser.Close()
			return
		case message := <-in:
			out <- r.publiser.Publish(
				[]byte(message),
				[]string{"my_routing_key"},
				rabbitmq.WithPublishOptionsContentType("application/json"),
				rabbitmq.WithPublishOptionsExchange("events"),
			)
		}
	}
}

func (r *RabbitMQ) consume(callback func(io.Reader) error, out chan<- error) {
	go func() {
		<-r.ctx.Done()
		r.consumer.Close()
	}()

	out <- r.consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		out <- callback(bytes.NewReader(d.Body))
		return rabbitmq.Ack
	})
}

func (r *RabbitMQ) Close() {
	r.cancel()
}
