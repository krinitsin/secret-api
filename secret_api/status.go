package secret_api

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"github.com/sirupsen/logrus"


	rpccodes "google.golang.org/grpc/codes"
)

// TODO: internal request cache to prevent multiple identical polling requests

func StatusByMac(remoteAddr, password, mac string) (*StatusResp, grpc.Error) {
	var err error
	if err = utils.ValidateIPv4HostPort(remoteAddr); err != nil {
		return nil, grpc.NewError(rpccodes.InvalidArgument, fmt.Errorf("remote_addr %s is invalid", remoteAddr))
	}
	if _, err = net.ParseMAC(remoteAddr); err != nil {
		return nil, grpc.NewError(rpccodes.InvalidArgument, fmt.Errorf("mac %s is invalid", mac))
	}

	req := mac2ipRequest(remoteAddr, password, mac)
	return pollStatus(req)
}

func StatusByIP(remoteAddr, password, userIP string) (*StatusResp, grpc.Error) {
	var err error
	if err = utils.ValidateIPv4HostPort(remoteAddr); err != nil {
		return nil, grpc.NewError(rpccodes.InvalidArgument, fmt.Errorf("remote_addr %s is invalid", remoteAddr))
	}
	if err = utils.ValidateIPv4HostPort(userIP); err != nil {
		return nil, grpc.NewError(rpccodes.InvalidArgument, fmt.Errorf("mac %s is invalid", userIP))
	}

	req := ip2macRequest(remoteAddr, password, userIP)
	return pollStatus(req)
}

// pollStatus отправляет запрос до тех пор пока не получит валидный результит или не наступит таймаут
func pollStatus(req *Request) (resp *StatusResp, err grpc.Error) {
	if delay := options.PollingDelay; delay != 0 {
		<-time.After(delay)
	}

	timeout := time.NewTimer(options.Timeout)
	ticker := time.NewTicker(options.PolligInterval)
	defer timeout.Stop()
	defer ticker.Stop()

	// отправляем запросы пока не получим валидный ответ
	for {
		tm := timeutils.NewTimeMeter()
		resp, err = getStatus(req)
		if err != nil || resp.UserIP == "" || resp.MAC == "" {
			metrics.AddSecretApiCall(req.method, "error", resp.Code, tm.Check())
			// на второй круг
			select {
			case <-timeout.C:
				err.UpdateCode(rpccodes.DeadlineExceeded)
				return nil, err
			case <-ticker.C:
			}
			continue
		}
		metrics.AddSecretApiCall(req.method, "success", resp.Code, tm.Check())
		return
	}
}

func getStatus(req *Request) (*StatusResp, grpc.Error) {
	body, err := req.send(apiRoute)
	if err != nil {
		return nil, grpc.NewError(rpccodes.Internal, err)
	}
	resp := &StatusResp{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, grpc.NewError(rpccodes.DeadlineExceeded, fmt.Errorf("unable to unmarshal response body: %s", err))
	}
	if err := errorByCode(resp.Code, resp.ReplyMessage); err != nil {
		return nil, err
	}
	logrus.
		WithField("type", req.Type).
		WithField("method", "POST").
		WithField("remote_ip", req.destAddr).
		WithField("user_ip", resp.UserIP).
		WithField("mac", resp.MAC).
		WithField("status_code", resp.Code).
		Debugf("Status request got reply: %s", resp.ReplyMessage)
	return resp, nil
}
