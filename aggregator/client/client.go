package client

import (
	"context"

	"github.com/unsuman/go-microservices/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, int64) (*types.Invoice, error)
}
