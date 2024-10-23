package main

import "github.com/unsuman/go-microservices/types"

type Aggregator interface {
	AggregateDistance(types.Distance) error
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
