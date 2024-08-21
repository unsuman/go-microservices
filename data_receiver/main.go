package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/unsuman/go-microservices/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
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
		dr.msgch <- data
	}
}
