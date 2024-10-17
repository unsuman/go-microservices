package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/types"
)

type DataProducer interface {
	ProduceData(*types.OBUData) error
}

type KafkaProducer struct {
	topic     string
	kafkaConn *kafka.Producer
}

var KafkaTopic string = "obu_data"

func NewKafkaProducer(topic string) (DataProducer, error) {

	conn, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
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

	err = kp.kafkaConn.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &KafkaTopic, Partition: kafka.PartitionAny},
		Value:          jsonData,
	}, nil)
	if err != nil {
		logrus.Fatal("failed to write messages:", err)
	}
	return nil
}
