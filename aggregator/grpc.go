package main

import (
	"context"

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

	return &types.None{}, g.svc.AggregateDistance(distance)
}

func (g *GRPCAggregator) GetInvoice(ctx context.Context, req *types.InvoiceRequest) (*types.InvoiceResponse, error) {
	invoice, err := g.svc.CalculateInvoice(req.ObuID)
	if err != nil {
		return nil, err
	}

	return &types.InvoiceResponse{
		ObuID:         int64(invoice.OBUid),
		TotalAmount:   invoice.TotalAmount,
		TotalDistance: invoice.TotalDistance,
	}, nil
}
