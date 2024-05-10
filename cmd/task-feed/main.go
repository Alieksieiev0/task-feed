package main

import (
	"log"
	"time"

	"github.com/Alieksieiev0/task-feed/internal/app"
	"github.com/Alieksieiev0/task-feed/internal/bot"
	"github.com/Alieksieiev0/task-feed/internal/broker"
	"github.com/Alieksieiev0/task-feed/internal/feed"
	"github.com/Alieksieiev0/task-feed/internal/model"
	"github.com/Alieksieiev0/task-feed/internal/repository"
	"github.com/Alieksieiev0/task-feed/internal/storage"
	"github.com/Alieksieiev0/task-feed/internal/streamer"
	"github.com/Alieksieiev0/task-feed/internal/transport"
)

func main() {
	db, err := storage.NewPostgreSQL().Open(model.Message{})
	if err != nil {
		log.Fatal(err)
	}
	feed := feed.NewJsonFeed(repository.NewGormRepository[model.Message](db))

	streamer := streamer.NewStringStreamer()
	broker, err := broker.NewRabbitMQBroker(
		storage.NewRabbitMQ(),
		"some",
		func(body []byte) error {
			message, err := feed.SaveMessage(body)
			if err != nil {
				return err
			}
			go streamer.AddMessage(message.String())
			return nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	server := transport.NewHttpServer(":3000", streamer, broker, feed)
	bot := bot.NewDelayBasedBot(time.Second*5, broker)

	twitterFeed := app.NewTwitterFeed(
		streamer,
		broker,
		server,
		bot,
		model.NewMessageFactory().CreateWithContent("Bot created message!"),
	)

	log.Fatal(twitterFeed.Run())
}
