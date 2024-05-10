package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	broker "github.com/Alieksieiev0/task-feed/action"
	"github.com/Alieksieiev0/task-feed/model"
	"github.com/Alieksieiev0/task-feed/repository"
	"github.com/Alieksieiev0/task-feed/storage"
	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Fatal(run())
}

var messages []string
var mu sync.Mutex

func run() error {
	db, err := storage.NewPostgreSQL().Open(model.Message{})
	if err != nil {
		return err
	}
	rep := repository.NewGormRepository[model.Message](db)
	app := fiber.New()
	conn, err := storage.NewRabbitMQ().Open()
	if err != nil {
		return err
	}
	publisher, err := broker.NewRabbitMQPublisher(conn)
	if err != nil {
		return err
	}
	consumer, err := broker.NewRabbitMQConsumer(conn, "some", func(body []byte) error {
		m := &model.Message{}
		err = json.Unmarshal(body, m)
		if err != nil {
			return err
		}
		return rep.Save(context.Background(), m)
	})
	if err != nil {
		return err
	}

	consumerErr := make(chan error)
	go func() {
		err := consumer.Run(consumerErr)
		if err != nil {
			publisher.Close()
		}
	}()
	go func() {
		for err := range consumerErr {
			fmt.Println(err)
		}
	}()

	fmt.Println("----")
	app.Post("/messages", func(c *fiber.Ctx) error {
		return publisher.Run(c.Body())
	})

	fmt.Println("----")
	app.Get("/messages", func(c *fiber.Ctx) error {
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			for i := 0; i < 10; i++ {
				fmt.Fprintf(w, "this is a message number %d", i)

				// Do not forget flushing streamed data to the client.
				if err := w.Flush(); err != nil {
					return
				}
				time.Sleep(time.Second)
			}
		})
		return nil
	})

	return app.Listen(":3000")
}
