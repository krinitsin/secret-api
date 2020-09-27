package http

import (
	"errors"
	nethttp "net/http"

	"github.com/sirupsen/logrus"
)

var (
	// server глобальный сервер
	server *Server
	// log логер
	log *logrus.Entry
)

// ErrUninitialized ошибка неинициализированности клиента
var ErrUninitialized = errors.New("uninitialized")

// Init инициализирует глобальный сервер
func Init(o *Options, appInfo string) {
	log = utils.NewLogger("http_server")
	server = New(o, appInfo)
}

// Start запускает глобальный сервер
func Start() {
	if server == nil {
		log.WithError(ErrUninitialized).Fatal("Unable to start server")
	}

	if err := server.Start(); err != nil && err != nethttp.ErrServerClosed {
		log.WithError(err).Fatal("Unable to start server")
	} else if err == nethttp.ErrServerClosed {
		log.Warnf("Server closed")
	}
}

// Stop останавливает глобальный сервер
func Stop() {
	if server == nil {
		log.WithError(ErrUninitialized).Fatal("Unable to stop server")
	}

	server.Stop()
	return
}
