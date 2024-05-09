package app

import "context"

type Repository[T any] interface {
	Get(ctx context.Context) ([]T, error)
	Save(ctx context.Context, entity *T) error
}

type Storage[T any] interface {
	Open(options ...any) (T, error)
}
