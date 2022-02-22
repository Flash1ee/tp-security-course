package utils

import (
	"http-proxy/cfg"

	"github.com/sirupsen/logrus"
)

func GetLogger(cfg *cfg.Config) *logrus.Logger {
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}
