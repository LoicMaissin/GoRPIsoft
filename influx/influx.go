package influx

/*
  Sends the data to the InfluxDB database
*/

import (
	"GoRPIsoft/analyser"
	"log"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

// Variables for the database and tags
const myDB = "test"

var serialNumber = os.Getenv("SERIALNB")
var importance = os.Getenv("IMPORTANCE")
var location = os.Getenv("LOCATION")

// WriteDB Creates batches of points from the actuator responses
func WriteDB(channelResponse chan [38]byte, channelBatches chan client.BatchPoints) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  myDB,
		Precision: "ms",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"actuator": serialNumber, "location": location, "importance": importance}
	last := [38]byte{}
	lastTime := time.Now()
	count := 0
	for response := range channelResponse {
		if last == response {
			lastTime = time.Now()
		} else {
			fields := analyser.ActuatorInfo(last)
			pt, errPt := client.NewPoint("measures", tags, fields, lastTime)
			if errPt != nil {
				log.Fatal(errPt)
			}
			bp.AddPoint(pt)

			fields = analyser.ActuatorInfo(response)
			pt, errPt = client.NewPoint("measures", tags, fields, time.Now())
			if errPt != nil {
				log.Fatal(errPt)
			}
			bp.AddPoint(pt)
			log.Println("Wrote a point!")

			last = response
			count++
		}
		if count == 200 {
			// Write the batch
			channelBatches <- bp
			bp, err = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  myDB,
				Precision: "ms",
			})
			if err != nil {
				log.Fatal(err)
			}
			count = 0
		}
	}
	fields := analyser.ActuatorInfo(last)
	pt, errPt := client.NewPoint("measures", tags, fields, lastTime)
	if errPt != nil {
		log.Fatal(errPt)
	}
	bp.AddPoint(pt)
	channelBatches <- bp
	time.Sleep(time.Second * 2)

	close(channelBatches)
}

// SendDB sends a batch of points to the database
func SendDB(channelBatches chan client.BatchPoints) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               "https://localhost:8086",
		Username:           os.Getenv("INFLUX_USER"),
		Password:           os.Getenv("INFLUX_PSSWD"),
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Println("Error creating client")
		log.Fatal(err)
	}
	defer c.Close() //ensure c is closed after function return
	for batch := range channelBatches {
		// Write the batch
		err := c.Write(batch)
		for err != nil {
			log.Println("Error client.Write()")
			log.Println(err)
			time.Sleep(time.Second / 2)
			err = c.Write(batch)
		}
		log.Println("Sent the batch!")
	}
}
