package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

)

const (
	ip2mac = "ip_to_mac"
	mac2ip = "mac_to_ip"
)

// Options опции сервера
type Options struct {
	Addr string
}

// Server реализация gRPC-сервера
type Server struct {
	options *Options
	gRPC    *grpc.Server
}

// New создает сервер
func New(o *Options) *Server {
	srv := new(Server)
	srv.options = o
	srv.gRPC = grpc.NewServer()
	RegisterSecretApiServer(srv.gRPC, srv)

	return srv
}

// Start запускает сервер
func (srv *Server) Start() error {
	listener, err := net.Listen("tcp", srv.options.Addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	return srv.gRPC.Serve(listener)
}

// Stop останавливает сервер
func (srv *Server) Stop() {
	srv.gRPC.GracefulStop()
}

// Ip2mac получает mac по ip
func (srv *Server) Ip2Mac(ctx context.Context, req *Ip2Mac_Request) ( *Response, error) {
	tm := timeutils.NewTimeMeter()
	log := utils.NewLogger(utils.UID())
	log.WithField("module", "ip2mac").WithField("remote_ip", req.RemoteAddr).WithField("user_ip", req.UserIp).Info("Request")

	statusResp, err := secret_api.StatusByIP(req.RemoteAddr, req.Password, req.UserIp)
	if err != nil {
		metrics.AddRequest("ip_to_mac", "grpc", "error", tm.Check())
		log.WithError(err).Error()
		return nil, err.ErrorWithCode()
	}

	metrics.AddRequest("ip_to_mac", "grpc", "success", tm.Check())
	log.WithField("method", ip2mac).Info("Mac info successfully acquired")
	return &Response{Mac: statusResp.MAC, Ip: statusResp.UserIP}, nil
}

func (srv *Server) Mac2Ip(ctx context.Context, req *Mac2Ip_Request) ( *Response, error) {
	tm := timeutils.NewTimeMeter()
	log := utils.NewLogger(utils.UID())
	log.WithField("module", "mac2ip").WithField("remote_ip", req.RemoteAddr).WithField("mac", req.Mac).Info("Request")

	statusResp, err := secret_api.StatusByMac(req.RemoteAddr, req.Password, req.Mac)
	if err != nil {
		metrics.AddRequest("mac_to_ip", "grpc", "error", tm.Check())
		log.WithError(err).Error()
		return nil, err.ErrorWithCode()
	}

	metrics.AddRequest("mac_to_ip", "grpc", "success", tm.Check())
	log.WithField("method", mac2ip).Info("IP info successfully acquired")
	return &Response{Mac: statusResp.MAC, Ip: statusResp.UserIP}, nil
}
