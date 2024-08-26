package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/unsuman/go-microservices/types"
)

type DataProducer interface {
	ProduceData(*types.OBUData) error
}

type KafkaProducer struct {
	topic     string
	kafkaConn *kafka.Conn
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	partition := 0
	address := fmt.Sprintf("%s:9092", os.Getenv("KAFKA_DOCKER_PORT"))
	conn, err := kafka.DialLeader(context.Background(), "tcp", address, kafkaTopic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	return &KafkaProducer{
		topic:     topic,
		kafkaConn: conn,
	}, nil
}

func (kp KafkaProducer) ProduceData(data *types.OBUData) error {
	_, err := kp.kafkaConn.WriteMessages(
		kafka.Message{WriterData: data},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	return nil
}
