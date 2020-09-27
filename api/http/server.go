package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

)

const (
	ip2mac = "ip_to_mac"
	mac2ip = "mac_to_ip"
)

// Options опции сервера
type Options struct {
	Addr string
}

// Server реализация HTTP-сервера
type Server struct {
	router  *mux.Router
	httpsrv *http.Server
	appInfo string
}

// New создает сервер
func New(o *Options, appInfo string) *Server {
	router := mux.NewRouter()

	httpsrv := &http.Server{
		Handler: router,
		Addr:    o.Addr,
	}

	s := &Server{
		router:  router,
		httpsrv: httpsrv,
		appInfo: appInfo,
	}

	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	router.HandleFunc("/version", s.versionHandler).Methods("GET")
	router.HandleFunc("/log/level", s.getLogLevelHandler).Methods("GET")
	router.HandleFunc("/log/level/{level}", s.setLogLevelHandler).Methods("POST")
	router.HandleFunc("/v1/secret/ip2mac", s.ip2macHandler).Methods("GET")
	router.HandleFunc("/v1/secret/mac2ip", s.mac2ipHandler).Methods("GET")

	return s
}

// Start запускает сервер
func (s *Server) Start() error {
	return s.httpsrv.ListenAndServe()
}

// Stop останавливает сервер
func (s *Server) Stop() error {
	return s.httpsrv.Shutdown(context.Background())
}

// versionHandler возвращает информацию о приложении
func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	log := utils.NewLogger(utils.UID())
	rw := httphelper.New(log, w, r)
	rw.ResponseString(http.StatusOK, s.appInfo)
}

// getLogHandler запрашивает уровень логирования
func (s *Server) getLogLevelHandler(w http.ResponseWriter, r *http.Request) {
	log := utils.NewLogger(utils.UID())
	rw := httphelper.New(log, w, r)
	rw.ResponseString(http.StatusOK, logrus.GetLevel().String())
}

// setLogHandler устанавливает уровень логирования
func (s *Server) setLogLevelHandler(w http.ResponseWriter, r *http.Request) {
	log := utils.NewLogger(utils.UID())
	rw := httphelper.New(log, w, r)

	vars := mux.Vars(r)
	level, ok := vars["level"]
	if !ok {
		rw.Response(http.StatusBadRequest, nil, nil)
		return
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		rw.ResponseError(http.StatusBadRequest, err)
		return
	}

	logrus.SetLevel(lvl)
	rw.Response(http.StatusNoContent, nil, nil)
}
