package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/KseniiaSalmina/Profiles/internal/api"
	"github.com/KseniiaSalmina/Profiles/internal/config"
	"github.com/KseniiaSalmina/Profiles/internal/database"
	"github.com/KseniiaSalmina/Profiles/internal/formatter"
	"github.com/KseniiaSalmina/Profiles/internal/logger"
	"github.com/sirupsen/logrus"
)

type Application struct {
	cfg       config.Application
	db        *database.Database
	formatter *formatter.Formatter
	logger    *logrus.Logger
	server    *api.Server
	closeCh   chan os.Signal
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
	if err := a.initDatabase(); err != nil {
		return err
	}

	a.initFormatter()
	if err := a.initLogger(); err != nil {
		return err
	}

	a.initServer()

	return nil
}

func (a *Application) initDatabase() error {
	db, err := database.NewDatabase(a.cfg.Database, a.cfg.Formatter.Salt)
	if err != nil {
		return fmt.Errorf("failed to init database")
	}

	a.db = db
	return nil
}

func (a *Application) initFormatter() {
	a.formatter = formatter.NewFormatter(a.cfg.Formatter, a.db)
}

func (a *Application) initLogger() error {
	l, err := logger.NewLogger(a.cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to init logger")
	}

	a.logger = l

	return nil
}

func (a *Application) initServer() {
	a.server = api.NewServer(a.cfg.Server, a.formatter, a.logger)
}

func (a *Application) Run() {
	defer a.stop()

	a.server.Run()

	<-a.closeCh
}

func (a *Application) stop() {
	if err := a.server.Shutdown(); err != nil {
		a.logger.Infof("server stopped: %s", err.Error())
	}
}

func (a *Application) readyToShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	a.closeCh = ch
}
