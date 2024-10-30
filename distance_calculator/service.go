package main

import (
	"math"
	"time"

	"github.com/unsuman/go-microservices/aggregator/client"
	"github.com/unsuman/go-microservices/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculatorService struct {
	points [][]float64
	client *client.HTTPClient
}

func NewCalculatorService(endPoint string) CalculatorServicer {
	return &CalculatorService{
		points: make([][]float64, 0),
		client: client.NewHTTPClient(endPoint),
	}
}

func (c *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0

	if len(c.points) > 0 {
		previousPoint := c.points[len(c.points)-1]
		distance = calculateDistance(previousPoint[0], data.Lat, previousPoint[1], data.Long)
	}
	c.points = append(c.points, []float64{data.Lat, data.Long})

	aggDistance := types.Distance{
		OBUID: data.OBUid,
		Value: distance,
		Unix:  time.Now().Unix(),
	}
	c.client.AggregateDistance(aggDistance)
	return distance, nil
}

func calculateDistance(x1, x2, y1, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
