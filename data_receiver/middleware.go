package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/types"
)

type LoggingMiddleware struct {
	prod DataProducer
}

func NewLoggingMiddleware(p DataProducer) *LoggingMiddleware {
	return &LoggingMiddleware{
		prod: p,
	}
}

func (l LoggingMiddleware) ProduceData(data *types.OBUData) error {
	t := time.Now()
	defer func(t time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUid,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(t),
		}).Info("producing to kafka")
	}(t)
	return l.prod.ProduceData(data)
}
