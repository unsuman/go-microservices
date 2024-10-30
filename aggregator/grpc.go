package main

import (
	"context"
	"fmt"

	"github.com/unsuman/go-microservices/types"
)

type GRPCAggregator struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregator(svc Aggregator) *GRPCAggregator {
	return &GRPCAggregator{
		svc: svc,
	}
}

func (g *GRPCAggregator) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int64(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}

	fmt.Println("Aggregating distance: ", distance)
	g.svc.AggregateDistance(distance)
	return &types.None{}, g.svc.AggregateDistance(distance)
}
