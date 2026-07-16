package grpcDelivery

import (
	"net"

	pb "github.com/ganfay/split-proto/pb"
	"google.golang.org/grpc"
)

type Server struct {
	addr       string
	grpcServer *grpc.Server
}

func NewServer(addr string) *Server {
	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, NewHandler())
	return &Server{addr: addr, grpcServer: s}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
