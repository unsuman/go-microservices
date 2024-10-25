package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/types"
)

type LoggingMiddleware struct {
	next Aggregator
}

func NewLoggingMiddleware(next Aggregator) Aggregator {
	return &LoggingMiddleware{
		next: next,
	}
}

func (l *LoggingMiddleware) AggregateDistance(d types.Distance) error {
	t := time.Now()
	err := l.next.AggregateDistance(d)
	defer func(t time.Time) {
		logrus.WithFields(logrus.Fields{
			"value": d.Value,
			"unix":  d.Unix,
			"obuID": d.OBUID,
			"error": err,
			"took":  time.Since(t),
		}).Info("aggregating distance")
	}(t)
	return nil
}

func (l *LoggingMiddleware) CalculateInvoice(obuId int64) (*types.Invoice, error) {
	t := time.Now()
	invoice, err := l.next.CalculateInvoice(obuId)
	if err != nil {
		return nil, err
	}
	defer func(t time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuid":         obuId,
			"totalAmount":   invoice.TotalAmount,
			"totalDistance": invoice.TotalDistance,
			"error":         err,
			"took":          time.Since(t),
		}).Info("Calculating invoice")
	}(t)
	return invoice, nil
}
