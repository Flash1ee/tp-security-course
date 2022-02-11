package utils

import "github.com/sirupsen/logrus"

func GetLogger() *logrus.Logger {
	log := logrus.New()
	return log
}
