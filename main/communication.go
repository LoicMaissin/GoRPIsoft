package main

import (
	"time"

	"GoRPIsoft/influx"
	"GoRPIsoft/serialBT"

	"github.com/influxdata/influxdb/client/v2"
)

func main() {

	channelResponse := make(chan [38]byte)
	channelBatches := make(chan client.BatchPoints)

	// Launch the Database client
	go influx.WriteDB(channelResponse, channelBatches)
	go influx.SendDB(channelBatches)

	for {
		channelResponse <- serialBT.GetAll()
		// You can use a time.Sleep() here to slow down the requests.
	}
	close(channelResponse)
	time.Sleep(time.Second * 3)
}
