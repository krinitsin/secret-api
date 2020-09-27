package utils

import (
	"errors"
	"net"
	"strconv"
)

func ValidateIPv4HostPort(addr string) (err error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if ip := net.ParseIP(host); ip == nil {
		return errors.New("bad ip format")
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return errors.New("bad port format")
	}
	if p > 65535 || p < 1 {
		return errors.New("port out of range")
	}
	return
}
