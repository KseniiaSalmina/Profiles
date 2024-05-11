package logger

import (
	"fmt"

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

	return l, nil
}
