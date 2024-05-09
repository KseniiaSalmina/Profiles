package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KseniiaSalmina/Profiles/internal/api"
	"github.com/KseniiaSalmina/Profiles/internal/config"
	"github.com/KseniiaSalmina/Profiles/internal/database"
)

type Application struct {
	cfg     config.Application
	db      *database.Database
	server  *api.Server
	closeCh chan os.Signal
}

func NewApplication(cfg config.Application) (*Application, error) {
	app := Application{
		cfg: cfg,
	}

	if err := app.bootstrap(); err != nil {
		return nil, err
	}

	app.readyToShutdown()

	return &app, nil
}

func (a *Application) bootstrap() error {
	//init dependencies
	if err := a.initDatabase(); err != nil {
		return err
	}

	//init controllers
	a.initServer()

	return nil
}

func (a *Application) initDatabase() error {
	db, err := database.NewDatabase(a.cfg.Database)
	if err != nil {
		return err
	}

	a.db = db
	return nil
}

func (a *Application) initServer() {
	a.server = api.NewServer(a.cfg.Server, a.db)
}

func (a *Application) Run() {
	defer a.stop()

	a.server.Run()

	<-a.closeCh
}

func (a *Application) stop() {
	if err := a.server.Shutdown(); err != nil {
		log.Println(err)
	}
}

func (a *Application) readyToShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	a.closeCh = ch
}
