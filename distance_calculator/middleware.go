package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/types"
)

type LoggingMiddleware struct {
	next CalculatorServicer
}

func NewLoggingMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LoggingMiddleware{
		next: next,
	}
}

func (l *LoggingMiddleware) CalculateDistance(data types.OBUData) (distance float64, err error) {
	t := time.Now()
	distance, err = l.next.CalculateDistance(data)
	defer func(t time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID":    data.OBUid,
			"error":    err,
			"distance": distance,
			"took":     time.Since(t),
		}).Info("calculating distance")
	}(t)
	return
}
