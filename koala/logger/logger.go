package logger

import (
	log "github.com/sirupsen/logrus"
)

var Logger *log.Entry

func init() {
	Logger = log.WithFields(log.Fields{})
}
