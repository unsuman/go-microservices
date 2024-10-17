package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/types"
)

type kafkaConsumer struct {
	reader *kafka.Reader
	topic  string
	calc   CalculatorServicer
}

func NewKafkaConsumer(topic string, s CalculatorServicer) *kafkaConsumer {
	address := fmt.Sprintf("%s:9092", os.Getenv("KAFKA_DOCKER_PORT"))
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     "obu_data", // Replace with your actual topic
		Partition: 0,          // Ensure you are reading from the correct partition
		MinBytes:  10e3,       // 10KB
		MaxBytes:  10e6,       // 10MB
	})

	return &kafkaConsumer{
		reader: r,
		topic:  topic,
		calc:   s,
	}
}

func (kc *kafkaConsumer) ConsumeData() {
	for {
		msg, err := kc.reader.ReadMessage(context.Background())
		if err != nil {
			logrus.Fatal("failed to unmarshal message:", err)
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Fatal("failed to unmarshal message:", err)
		}

		fmt.Printf("Received OBUData: %+v\n", data)

		_, err = kc.calc.CalculateDistance(data)
		if err != nil {
			logrus.Fatal("failed to calculate distance: ", err)
		}
	}
}
