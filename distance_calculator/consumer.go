package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/types"
)

type kafkaConsumer struct {
	reader *kafka.Consumer
	topic  string
	calc   CalculatorServicer
}

func NewKafkaConsumer(topic string, s CalculatorServicer) *kafkaConsumer {
	r, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logrus.Fatal("failed to dial leader:", err)
	}

	err = r.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logrus.Fatal("failed to subscribe to topic:", err)
	}

	return &kafkaConsumer{
		reader: r,
		topic:  topic,
		calc:   s,
	}
}

func (kc *kafkaConsumer) ConsumeData() {
	for {
		msg, err := kc.reader.ReadMessage(-1)
		if err != nil {
			logrus.Fatal("failed to consume message:", err)
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
