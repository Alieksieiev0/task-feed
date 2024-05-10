package transport

import (
	"bufio"
	"context"

	"github.com/Alieksieiev0/task-feed/internal/app"
	"github.com/Alieksieiev0/task-feed/internal/model"
	"github.com/gofiber/fiber/v2"
)

func NewHttpServer(
	addr string,
	streamer app.Streamer[string],
	broker app.Broker[[]byte],
	feed app.Feed[[]byte, model.Message],
) *HttpServer {
	return &HttpServer{
		app:      fiber.New(),
		addr:     addr,
		streamer: streamer,
		broker:   broker,
		feed:     feed,
	}
}

type HttpServer struct {
	app      *fiber.App
	addr     string
	streamer app.Streamer[string]
	broker   app.Broker[[]byte]
	feed     app.Feed[[]byte, model.Message]
}

func (h *HttpServer) Run() error {
	h.app.Post("/messages", func(c *fiber.Ctx) error {
		return h.broker.Publish(c.Body())
	})

	h.app.Get("/messages", func(c *fiber.Ctx) error {
		messages, err := h.feed.GetMessages()
		if err != nil {
			return err
		}
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			ctx, cancel := context.WithCancel(context.Background())
			client := NewHttpClient(w, cancel)
			for _, m := range messages {
				if err := client.Write(m.String()); err != nil {
					return
				}
			}
			h.streamer.AddClient(client)
			for range ctx.Done() {
				return
			}
		})
		return nil
	})

	return h.app.Listen(h.addr)
}

func (h *HttpServer) Close() error {
	return h.app.Shutdown()
}
