package grpcserver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	_defaultNet  = "tcp"
	_defaultAddr = ":8000"
)

var RpcServer *grpc.Server

// Server -.
type Server struct {
	server *grpc.Server
	notify chan error
	Addr   string
}

func New(opts ...Option) *Server {
	grpcServer := grpc.NewServer()

	// 注册 grpcurl 所需的 reflection 服务
	reflection.Register(grpcServer)

	RpcServer = grpcServer

	s := &Server{
		server: grpcServer,
		notify: make(chan error, 1),
		Addr:   _defaultAddr,
	}

	for _, opt := range opts {
		opt(s)
	}
	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		listen, err := net.Listen(_defaultNet, s.Addr)
		if err != nil {
			log.Fatal("net.Listen err: %v", err)
		}
		s.notify <- s.server.Serve(listen)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() {
	s.server.GracefulStop()
}
