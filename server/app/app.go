package app

import (
	"github.com/mojtabamovahedi/chatroom/server/config"
	chatMap "github.com/mojtabamovahedi/chatroom/server/pkg/map"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nats"
)

type App struct {
	cfg         config.Config
	natsClient  nats.Nats
	userMap     *chatMap.Map[string, string]
	chatroomMap *chatMap.Map[string, string]
}

func NewApp(cfg config.Config) (*App, error) {
	app := &App{
		cfg: cfg,
	}

	err := app.setNats()
	if err != nil {
		return nil, err
	}

	app.userMap = chatMap.NewMap[string, string]()
	app.chatroomMap = chatMap.NewMap[string, string]()

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

func (app *App) MapUser() *chatMap.Map[string, string] {
	return app.userMap
}

func (app *App) MapChatroom() *chatMap.Map[string, string] {
	return app.chatroomMap
}

func (app *App) Nats() *nats.Nats {
	return &app.natsClient
}
