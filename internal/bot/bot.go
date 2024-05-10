package bot

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Alieksieiev0/task-feed/internal/app"
)

func NewDelayBasedBot(delay time.Duration, broker app.Broker[[]byte]) *DelayBasedBot {
	ctx, cancel := context.WithCancel(context.Background())
	return &DelayBasedBot{delay: delay, broker: broker, ctx: ctx, cancel: cancel}
}

type DelayBasedBot struct {
	delay  time.Duration
	broker app.Broker[[]byte]
	ctx    context.Context
	cancel context.CancelFunc
}

func (d *DelayBasedBot) Run(entity app.Model) error {
	for {
		select {
		case <-d.ctx.Done():
			return nil
		default:
			time.Sleep(d.delay)
			encoded, err := json.Marshal(&entity)
			if err != nil {
				return err
			}

			if err = d.broker.Publish(encoded); err != nil {
				return err
			}
		}
	}

}

func (d *DelayBasedBot) Close() {
	d.cancel()
}
