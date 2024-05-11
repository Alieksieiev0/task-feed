package app

import "context"

type Closeable interface {
	Close()
}

type ErrorProneCloseable interface {
	Close() error
}

type Client[T any] interface {
	Write(data T) error
	Cancel()
}

type Streamer[T any] interface {
	AddClient(client Client[T])
	AddMessage(msg T)
	Stream()
	Closeable
}

type Feed[Input any, Output any] interface {
	SaveMessage(data Input) (Output, error)
	GetMessages() ([]Output, error)
}

type Broker[T any] interface {
	Publish(input T) error
	Consume() error
	Closeable
}

type Server interface {
	Run(addr string) error
	ErrorProneCloseable
}

type Bot[T any] interface {
	Run(T) error
	Closeable
}

type Repository[T any, V any] interface {
	Get(ctx context.Context, opts ...V) ([]T, error)
	Save(ctx context.Context, entity *T, opts ...V) error
}

type Storage[Input any, Output any] interface {
	Open(options ...Input) (Output, error)
}

type Model interface {
	String() string
}
