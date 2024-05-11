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

func (e *GormRepository[T]) Get(
	ctx context.Context,
	opts ...func(db *gorm.DB) *gorm.DB,
) ([]T, error) {
	entities := []T{}
	return entities, e.applyParams(e.db, opts...).Find(&entities).Error
}

func (e *GormRepository[T]) Save(
	ctx context.Context,
	entity *T,
	opts ...func(db *gorm.DB) *gorm.DB,
) error {
	return e.applyParams(e.db, opts...).Save(entity).Error
}

func (e *GormRepository[T]) applyParams(db *gorm.DB, opts ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}
