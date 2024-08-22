package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	partition := 0

	//My Kafka is running on a different computer using Docker
	address := fmt.Sprintf("%s:9092", os.Getenv("KAFKA_DOCKER_PORT"))

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, kafkaTopic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	datarec := NewDataReceiver()
	http.HandleFunc("/ws", datarec.wsHandler)
	http.ListenAndServe(":30000", nil)

	defer datarec.conn.Close()
}

type DataReceiver struct {
	conn  *websocket.Conn
	msgch chan types.OBUData
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
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

func (dr DataReceiver) readWsReceiveloop() {
	fmt.Println("New OBU connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("--- received data from OBU [%d] :: lat[%.2f] long[%.2f] \n", data.OBUid, data.Lat, data.Long)
		// dr.msgch <- data
	}
}
