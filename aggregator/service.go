package main

import "github.com/unsuman/go-microservices/types"

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int64) (*types.Invoice, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(d types.Distance) error {
	return i.store.InsertDistance(d)
}

func (i *InvoiceAggregator) CalculateInvoice(obuID int64) (*types.Invoice, error) {
	dist, err := i.store.GetDistance(obuID)
	if err != nil {
		return nil, err
	}

	return &types.Invoice{
		OBUid:         obuID,
		TotalDistance: dist,
		TotalAmount:   dist * 1.5,
	}, nil

}
