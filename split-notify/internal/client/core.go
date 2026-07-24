package client

import (
	"context"
	"fmt"
	"log"

	pb "github.com/ganfay/split-proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CoreClient struct {
	client pb.NotificationServiceClient
	conn   *grpc.ClientConn
}

func NewCoreClient(target string) (*CoreClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("client: failed to connect to gRPC server: %w", err)
	}

	cl := pb.NewNotificationServiceClient(conn)

	return &CoreClient{client: cl, conn: conn}, nil
}

func (c *CoreClient) GetTargets(ctx context.Context, fundID int64, creatorID int64) ([]*pb.Target, error) {
	targets, err := c.client.GetNotificationTargets(ctx, &pb.GetRequest{FundId: fundID, CreatorIid: creatorID})
	if err != nil {
		return nil, err
	}
	log.Printf("Core: Get targets successfully on gRPC connection")
	listTargets := targets.GetTargets()
	log.Printf("Core: Targets: %v", listTargets)
	return listTargets, err
}

func (c *CoreClient) Ping(ctx context.Context, name string) (*pb.PingReply, error) {
	ping, err := c.client.SayPing(ctx, &pb.PingRequest{Name: name})
	return ping, err
}

func (c *CoreClient) NotificateTargets(ctx context.Context, targets []*pb.Target) error {
	for _, t := range targets {
		text := fmt.Sprintf("Hello, %s! *someone add new expense in *someone fund.", t.FirstName) //надо крейтор нейм добавтиь и фанд нейм
		_, err := c.client.NotificateTarget(ctx, &pb.NotTarget{Msg: text, TgId: t.TgId})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CoreClient) Close() error {
	return c.conn.Close()
}
