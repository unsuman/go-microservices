package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/unsuman/go-microservices/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const kafkaTopic = "obu-data"

func main() {
	datarec := NewDataReceiver()
	http.HandleFunc("/ws", datarec.wsHandler)
	http.ListenAndServe(":30000", nil)

	defer datarec.conn.Close()
}

type DataReceiver struct {
	conn      *websocket.Conn
	msgch     chan types.OBUData
	kafkaConn *kafka.Conn
}

func NewDataReceiver() *DataReceiver {
	partition := 0
	address := fmt.Sprintf("%s:9092", os.Getenv("KAFKA_DOCKER_PORT"))

	kafkaConn, err := kafka.DialLeader(context.Background(), "tcp", address, kafkaTopic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	return &DataReceiver{
		msgch:     make(chan types.OBUData, 128),
		kafkaConn: kafkaConn,
	}
}

func (dr DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.readWsReceiveloop()
}

func (dr DataReceiver) produceData(data *types.OBUData) error {
	// dr.kafkaConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := dr.kafkaConn.WriteMessages(
		kafka.Message{Value: []byte(string(data.OBUid))},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	// if err := dr.kafkaConn.Close(); err != nil {
	// 	log.Fatal("failed to close writer:", err)
	// }

	return nil
}

func (dr DataReceiver) readWsReceiveloop() {
	fmt.Println("New OBU connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("--- received data from OBU [%d] :: lat[%.2f] long[%.2f] \n", data.OBUid, data.Lat, data.Long)
		if err := dr.produceData(&data); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Data produced to Kafka")
	}
}
