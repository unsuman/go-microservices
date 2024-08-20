package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

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
	for {
		for i := 0; i < 25; i++ {
			lat, long := generateLatLong()
			data := OBUData{
				OBUid: generateOBUID(),
				Lat:   lat,
				Long:  long,
			}
			fmt.Println(data)
		}
		time.Sleep(sendInterval)
	}
}
