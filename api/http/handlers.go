package http

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

)

// ip2macHandler определение mac по ip
func (s *Server) ip2macHandler(w http.ResponseWriter, r *http.Request) {
	tm := timeutils.NewTimeMeter()
	log := utils.NewLogger(utils.UID())
	rw := httphelper.New(log, w, r)

	remoteAddr := rw.QueryValue("remote_addr")
	password := rw.QueryValue("password")
	userIP := rw.QueryValue("user_ip")

	log = log.WithField("request_type", ip2mac).WithField("api", "http").WithField("remote_addr", remoteAddr).WithField("user_ip", userIP)
	log.Info("Request")

	var err error
	defer func() {
		publishMetricsAndLog(ip2mac, tm.Check(), log, err)
	}()

	if err := utils.ValidateIPv4HostPort(remoteAddr); err != nil {
		rw.ResponseString(http.StatusBadRequest, fmt.Sprintf("remote_addr %s is invalid", remoteAddr))
		return
	}
	if err := utils.ValidateIPv4HostPort(userIP); err != nil {
		rw.ResponseString(http.StatusBadRequest, fmt.Sprintf("user_ip %s is invalid", userIP))
		return
	}

	resp, err := secret_api.StatusByIP(remoteAddr, password, userIP)
	if err != nil {
		rw.Response(http.StatusGone, nil, nil)
		return
	}
	mac := common.MACFormat(resp.MAC)
	log.WithField("mac", mac).Info("Mac successfully found")

	rsp := &Response{
		Userip: resp.UserIP,
		Mac:    mac,
	}

	rw.ResponseJSON(http.StatusOK, rsp)
}

// mac2ipHandler определение ip по mac
func (s *Server) mac2ipHandler(w http.ResponseWriter, r *http.Request) {
	tm := timeutils.NewTimeMeter()
	log := utils.NewLogger(utils.UID())
	rw := httphelper.New(log, w, r)

	remoteAddr := rw.QueryValue("remote_addr")
	password := rw.QueryValue("password")
	mac := rw.QueryValue("mac")

	log = log.WithField("request_type", mac2ip).WithField("api", "http").WithField("remote_addr", remoteAddr).WithField("mac", mac)
	log.Info("Request")

	var err error
	defer func() {
		publishMetricsAndLog(mac2ip, tm.Check(), log, err)
	}()

	if err := utils.ValidateIPv4HostPort(remoteAddr); err != nil {
		rw.ResponseString(http.StatusBadRequest, fmt.Sprintf("remote_addr %s is invalid", remoteAddr))
		return
	}
	if _, err := net.ParseMAC(mac); err != nil {
		rw.ResponseString(http.StatusBadRequest, fmt.Sprintf("mac %s is invalid", mac))
		return
	}

	resp, err := secret_api.StatusByMac(remoteAddr, password, mac)
	if err != nil {
		rw.Response(http.StatusGone, nil, nil)
		return
	}

	log.WithField("user_ip", resp.UserIP).Info("User_ip successfully found")

	rsp := &Response{
		Userip: resp.UserIP,
		Mac:    mac,
	}

	rw.ResponseJSON(http.StatusOK, rsp)
}

func publishMetricsAndLog(handler string, t time.Duration, log *logrus.Entry, err error) {
	status := "error"
	if err != nil {
		log.WithError(err).Error()
	}
	metrics.AddRequest(handler, "http", status, t)
}
