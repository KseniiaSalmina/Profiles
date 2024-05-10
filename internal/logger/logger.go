package logger

import (
	"fmt"
	"os"

	"github.com/KseniiaSalmina/Profiles/internal/config"
	"github.com/sirupsen/logrus"
)

func NewLogger(cfg config.Logger) (*logrus.Logger, error) {
	l := logrus.New()

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parce log level: %w", err)
	}

	l.SetLevel(lvl)

	file, err := os.OpenFile("profiles.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Info("Failed to log to file, using default stderr")
		return l, nil
	}

	l.SetOutput(file)

	return l, nil
}
