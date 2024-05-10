package streamer

import (
	"slices"

	"github.com/Alieksieiev0/task-feed/internal/app"
)

func NewStringStreamer[T string]() *StringStreamer[T] {
	return &StringStreamer[T]{message: make(chan T)}
}

type StringStreamer[T string] struct {
	clients []app.Client[T]
	message chan T
}

func (t *StringStreamer[T]) AddClient(client app.Client[T]) {
	t.clients = append(t.clients, client)
}

func (t *StringStreamer[T]) AddMessage(msg T) {
	t.message <- msg
}

func (t *StringStreamer[T]) Stream() {
	for m := range t.message {
		var invalidClients []int
		for i, c := range t.clients {
			if err := c.Write(m); err != nil {
				c.Cancel()
				invalidClients = append(invalidClients, i)
			}
		}

		for j, i := range invalidClients {
			delAt := i - j
			t.clients = slices.Delete(t.clients, delAt, delAt+1)
		}
	}
}

func (t *StringStreamer[T]) Close() {
	close(t.message)
}
