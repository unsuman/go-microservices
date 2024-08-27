package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/types"
)

func main() {
	address := fmt.Sprintf("%s:9092", os.Getenv("KAFKA_DOCKER_PORT"))

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     "obu_data", // Replace with your actual topic
		Partition: 0,          // Ensure you are reading from the correct partition
		MinBytes:  10e3,       // 10KB
		MaxBytes:  10e6,       // 10MB
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("failed to read messages:", err)
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			log.Fatal("failed to unmarshal message:", err)
		}

		fmt.Printf("Received OBUData: %+v\n", data)
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUid,
			"lat":   data.Lat,
			"long":  data.Long,
		}).Info("received from kafka")
	}
}
