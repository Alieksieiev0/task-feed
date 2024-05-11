package transport

import (
	"bufio"
	"context"

	"github.com/Alieksieiev0/task-feed/internal/app"
	"github.com/Alieksieiev0/task-feed/internal/model"
	"github.com/gofiber/fiber/v2"
)

func NewHttpServer(
	streamer app.Streamer[string],
	broker app.Broker[[]byte],
	feed app.Feed[[]byte, model.Message],
) *HttpServer {
	return &HttpServer{
		app:      fiber.New(),
		streamer: streamer,
		broker:   broker,
		feed:     feed,
	}
}

type HttpServer struct {
	app      *fiber.App
	streamer app.Streamer[string]
	broker   app.Broker[[]byte]
	feed     app.Feed[[]byte, model.Message]
}

func (h *HttpServer) Run(addr string) error {
	h.app.Post("/messages", func(c *fiber.Ctx) error {
		return h.broker.Publish(c.Body())
	})

	h.app.Get("/messages", func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		messages, err := h.feed.GetMessages()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
			<-ctx.Done()
			return
		})
		return nil
	})

	return h.app.Listen(addr)
}

func (h *HttpServer) Close() error {
	return h.app.Shutdown()
}
