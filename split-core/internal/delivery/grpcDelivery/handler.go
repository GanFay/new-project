package grpcDelivery

import (
	"context"
	"log/slog"

	"github.com/ganfay/split-core/internal/usecase"
	pb "github.com/ganfay/split-proto/pb"
)

type Handler struct {
	pb.UnimplementedNotificationServiceServer
	fundUC usecase.FundUsecase
}

func NewHandler(fundUC usecase.FundUsecase) *Handler {
	return &Handler{fundUC: fundUC}
}

func (h *Handler) SayPing(_ context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	slog.Info("Received: %v", in.GetName())
	return &pb.PingReply{Message: "Ping " + in.GetName()}, nil
}

func (h *Handler) GetNotificationTargets(ctx context.Context, in *pb.GetRequest) (*pb.TargetsResponse, error) {
	fundID := in.GetFundId()
	creatorIID := in.GetCreatorIid()
	slog.Info("Received: fund_id: ", fundID, "creator_id:", creatorIID)
	members, err := h.fundUC.GetMembers(ctx, int(fundID))
	if err != nil {
		slog.Error("Error to get members for notification service", "err", err)
		return nil, err
	}
	var targets []*pb.Target
	for _, m := range members {
		slog.Debug("TargetUser", "target", m)
		if m.ID == creatorIID || m.IsVirtual {
			continue
		}
		var tgID int64
		if m.TgID != nil {
			tgID = *m.TgID
		}

		target := &pb.Target{
			Iid:       m.ID,
			FirstName: m.FirstName,
			TgId:      tgID,
		}

		targets = append(targets, target)
	}
	return &pb.TargetsResponse{Targets: targets}, err
}
