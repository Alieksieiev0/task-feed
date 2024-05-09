package repository

import (
	"context"

	"gorm.io/gorm"
)

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

type GormRepository[T any] struct {
	db *gorm.DB
}

func (e *GormRepository[T]) Get(ctx context.Context) ([]T, error) {
	entities := []T{}
	return entities, e.db.Find(&entities).Error
}

func (e *GormRepository[T]) Save(ctx context.Context, entity *T) error {
	return e.db.Save(entity).Error
}
