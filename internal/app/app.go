package app

import "golang.org/x/sync/errgroup"

func NewTwitterFeed(
	streamer Streamer[string],
	broker Broker[[]byte],
	server Server,
	bot Bot[Model],
	templateModel Model,
) TwitterFeed {
	return TwitterFeed{
		streamer:      streamer,
		broker:        broker,
		server:        server,
		bot:           bot,
		templateModel: templateModel,
	}
}

type TwitterFeed struct {
	streamer      Streamer[string]
	broker        Broker[[]byte]
	server        Server
	bot           Bot[Model]
	templateModel Model
}

func (t TwitterFeed) Run() error {
	var g errgroup.Group

	g.Go(func() error {
		t.streamer.Stream()
		return nil
	})
	g.Go(func() error {
		return t.bot.Run(t.templateModel)
	})
	g.Go(t.broker.Consume)
	g.Go(t.server.Run)

	if err := g.Wait(); err != nil {
		t.streamer.Close()
		t.bot.Close()
		t.broker.Close()
		if servErr := t.server.Close(); servErr != nil {
			return servErr
		}
		return err
	}
	return nil
}
