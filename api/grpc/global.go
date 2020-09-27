package grpc

import (
	"errors"

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
func Init(o *Options) {
	log = utils.NewLogger("grpc_server")
	server = New(o)
}

// Start запускает глобальный сервер
func Start() {
	if server == nil {
		log.WithError(ErrUninitialized).Fatal("Unable to start server")
	}

	if err := server.Start(); err != nil {
		log.WithError(err).Fatal("Unable to start server")
	}
}

// Stop останавливает глобальный сервер
func Stop() {
	if server == nil {
		log.WithError(ErrUninitialized).Fatal("Unable to stop server")
	}

	server.Stop()
	log.Warn("Server closed")
	return
}
