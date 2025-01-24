package app

import (
	"github.com/mojtabamovahedi/chatroom/server/config"
	chatMap "github.com/mojtabamovahedi/chatroom/server/pkg/map"
	"github.com/mojtabamovahedi/chatroom/server/pkg/map/types"
	"github.com/mojtabamovahedi/chatroom/server/pkg/nats"
)

// App represents the application with its dependencies
type App struct {
	cfg         config.Config
	natsClient  *nats.Nats
	chatroomMap *chatMap.Map[string, *types.ChatRoom]
	userMap     *chatMap.Map[string, *types.User]
}

// NewApp creates a new application instance
func NewApp(cfg config.Config) (*App, error) {
	app := &App{
		cfg: cfg,
	}

	err := app.setNats()
	if err != nil {
		return nil, err
	}

	app.chatroomMap = chatMap.NewMap[string, *types.ChatRoom]()
	app.userMap = chatMap.NewMap[string, *types.User]()

	return app, nil
}

// MustNewApp creates a new application instance and panics if there is an error
func MustNewApp(cfg config.Config) *App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}

// setNats initializes the NATS client
func (app *App) setNats() error {
    n, err := nats.New(app.cfg.Nats.Host, app.cfg.Nats.Port)
    if err != nil {
        return err
    }
    app.natsClient = n
    return nil
}

// MapUser returns the user map
func (app *App) MapUser() *chatMap.Map[string, *types.User] {
    return app.userMap
}

// MapChatroom returns the chatroom map
func (app *App) MapChatroom() *chatMap.Map[string, *types.ChatRoom] {
    return app.chatroomMap
}

// Nats returns the NATS client
func (app *App) Nats() *nats.Nats {
    return app.natsClient
}

// Shutdown closes the NATS client connection
func (app *App) Shutdown() {
    if app.natsClient != nil {
        app.natsClient.Close()
    }
}