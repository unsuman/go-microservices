package client

import (
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
