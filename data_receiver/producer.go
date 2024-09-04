package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/types"
)

type DataProducer interface {
	ProduceData(*types.OBUData) error
}

type KafkaProducer struct {
	topic     string
	kafkaConn *kafka.Conn
}

const KafkaTopic = "obu_data"

func NewKafkaProducer(topic string) (DataProducer, error) {
	partition := 0
	address := fmt.Sprintf("%s:9092", os.Getenv("KAFKA_DOCKER_PORT"))
	conn, err := kafka.DialLeader(context.Background(), "tcp", address, KafkaTopic, partition)
	if err != nil {
		logrus.Fatal("failed to dial leader:", err)
	}
	return &KafkaProducer{
		topic:     topic,
		kafkaConn: conn,
	}, nil
}

func (kp KafkaProducer) ProduceData(data *types.OBUData) error {
	jsonData, err := json.Marshal(&data)
	if err != nil {
		logrus.Fatal("failed to marshal data:", err)
	}

	_, err = kp.kafkaConn.WriteMessages(
		kafka.Message{Value: jsonData},
	)
	if err != nil {
		logrus.Fatal("failed to write messages:", err)
	}
	return nil
}
