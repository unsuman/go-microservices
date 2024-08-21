package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/unsuman/go-microservices/types"
)

const wsEndpoint = "ws://127.0.0.1:30000/ws"
const sendInterval = time.Second

type OBUData struct {
	OBUid int64   `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

func generateLatLong() (float64, float64) {
	return generateCoord(), generateCoord()
}

func generateCoord() float64 {
	coord := float64(rand.Intn(100) + 1)
	return coord + rand.Float64()
}

func generateOBUID() int64 {
	return int64(rand.Intn(math.MaxInt))
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < 25; i++ {
			lat, long := generateLatLong()
			data := types.OBUData{
				OBUid: generateOBUID(),
				Lat:   lat,
				Long:  long,
			}
			time.Sleep(sendInterval)
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
