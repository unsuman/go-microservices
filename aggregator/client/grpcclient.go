package client

import (
	"context"
	"log"

	"github.com/unsuman/go-microservices/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRPCClient(endpoint string) *GrpcClient {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	c := types.NewAggregatorClient(conn)
	return &GrpcClient{
		Endpoint:         endpoint,
		AggregatorClient: c,
	}
}

func (c *GrpcClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	_, err := c.AggregatorClient.Aggregate(ctx, aggReq)
	if err != nil {
		return err
	}
	return nil
}

func (c *GrpcClient) GetInvoice(ctx context.Context, obuID int64) (*types.Invoice, error) {
	return &types.Invoice{
		OBUid:         obuID,
		TotalAmount:   1003.21,
		TotalDistance: 21.121,
	}, nil
}
