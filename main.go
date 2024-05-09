package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/Alieksieiev0/task-feed/broker"
	"github.com/Alieksieiev0/task-feed/model"
	"github.com/Alieksieiev0/task-feed/repository"
	"github.com/Alieksieiev0/task-feed/storage"
	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Fatal(run())
}

func run() error {
	db, err := storage.NewPostgreSQL().Open(model.Message{})
	if err != nil {
		return err
	}
	rep := repository.NewGormRepository[model.Message](db)
	app := fiber.New()
	brok, err := broker.NewRabbitMQ("some", "some")
	if err != nil {
		return err
	}
	msgChan, errChan := brok.Run(func(r io.Reader) error {
		m := &model.Message{}
		err := json.NewDecoder(r).Decode(m)
		if err != nil {
			return err
		}
		rep.Save(context.Background(), m)
		return nil
	})

	go func() {
		for err := range errChan {
			fmt.Println(err)
		}
	}()

	app.Post("/messages", func(c *fiber.Ctx) error {
		go func() {
			msgChan <- string(c.Body())
		}()
		return nil
	})

	return app.Listen(":3000")
}
