package storage

import (
	"fmt"
	"os"

	"github.com/wagslane/go-rabbitmq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgreSQL() *PostgreSQL {
	return &PostgreSQL{}
}

type PostgreSQL struct {
}

func (p *PostgreSQL) Open(entities ...any) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return db, db.AutoMigrate(entities...)
}

func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{}
}

type RabbitMQ struct {
}

func (r *RabbitMQ) Open() (*rabbitmq.Conn, error) {
	return rabbitmq.NewConn(
		fmt.Sprintf(
			"amqp://%s:%s@%s",
			os.Getenv("RABBITMQ_USER"),
			os.Getenv("RABBITMQ_PASS"),
			os.Getenv("RABBITMQ_HOST"),
		),
		rabbitmq.WithConnectionOptionsLogging,
	)
}
