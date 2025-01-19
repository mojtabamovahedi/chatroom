package app

import (
	"github.com/mojtabamovahedi/chatroom/server/config"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nats"
)

type App struct {
	cfg        config.Config
	natsClient nats.Nats
}

func NewApp(cfg config.Config) (*App, error) {
	app := &App{
		cfg: cfg,
	}

	err := app.setNats()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func MustNewApp(cfg config.Config) *App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}

func (app *App) setNats() error {
	n, err := nats.New(app.cfg.Nats.Host, app.cfg.Nats.Port)
	if err != nil {
		return err
	}
	app.natsClient = *n
	return nil
}
