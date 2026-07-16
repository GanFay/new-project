package grpcDelivery

import (
	"context"
	"log/slog"

	pb "github.com/ganfay/split-proto/pb"
)

type Handler struct {
	pb.UnimplementedNotificationServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SayPing(_ context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	slog.Info("Received: %v", in.GetName())
	return &pb.PingReply{Message: "Ping " + in.GetName()}, nil
}
