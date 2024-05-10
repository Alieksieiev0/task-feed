package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        string         `gorm:"type:uuid" json:"id"`
	CreatedAt time.Time      `                 json:"created_at"`
	UpdatedAt time.Time      `                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"     json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}

type Message struct {
	Base
	Content string `json:"content" gorm:"default:null;not null;"`
}

func (m Message) String() string {
	return fmt.Sprintf("[%s]: %s", m.CreatedAt, m.Content)
}
