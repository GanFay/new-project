package processor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ganfay/split-notify/internal/client"
)

type Processor struct {
	grpcClient *client.CoreClient
}

type ExpenseCreatedEvent struct {
	FundID    int64 `json:"fund_id"`
	CreatorID int64 `json:"creator_id"`
}

func NewProcessor(grpcClient *client.CoreClient) *Processor {
	return &Processor{grpcClient: grpcClient}
}

func (p *Processor) HandleMessage(ctx context.Context, body []byte) error {
	log.Printf("Processor: starting handle event... %s", body)

	var event ExpenseCreatedEvent

	err := json.Unmarshal(body, &event)
	if err != nil {
		return err
	}
	log.Printf("Processor: Get event! Fund ID: %d, Creator ID: %d", event.FundID, event.CreatorID)
	return nil
}
