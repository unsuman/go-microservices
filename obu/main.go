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
const sendInterval = time.Second * 5

func generateLatLong() (float64, float64) {
	return generateCoord(), generateCoord()
}

func generateCoord() float64 {
	coord := float64(rand.Intn(100) + 1)
	return coord + rand.Float64()
}

func generateOBUID(n int) []int64 {
	ids := make([]int64, n)
	for i := 0; i < n; i++ {
		ids[i] = int64(rand.Intn(math.MaxInt))
	}
	return ids
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	ids := generateOBUID(5)
	for {
		for _, id := range ids {
			lat, long := generateLatLong()
			data := types.OBUData{
				OBUid: id,
				Lat:   lat,
				Long:  long,
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
			time.Sleep(sendInterval)
		}
	}
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
