package feed

import (
	"context"
	"encoding/json"

	"github.com/Alieksieiev0/task-feed/internal/app"
	"gorm.io/gorm"
)

func NewJsonFeed[T app.Model](repository app.Repository[T, func(*gorm.DB) *gorm.DB]) *JsonFeed[T] {
	return &JsonFeed[T]{repository: repository}
}

type JsonFeed[T app.Model] struct {
	repository app.Repository[T, func(*gorm.DB) *gorm.DB]
}

func (j *JsonFeed[T]) SaveMessage(data []byte) (T, error) {
	m := new(T)
	err := json.Unmarshal(data, m)
	if err != nil {
		return *m, err
	}

	return *m, j.repository.Save(context.Background(), m)
}

func (j *JsonFeed[T]) GetMessages() ([]T, error) {
	return j.repository.Get(context.Background(), func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at")
	})
}
