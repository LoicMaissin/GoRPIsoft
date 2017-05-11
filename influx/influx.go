package influx

/*
  Sends the data to the InfluxDB database
*/

import (
	"GoRPIsoft/analyser"
	"bufio"
	"log"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var factory, importance, area, serialNumber, tagNumber, model, name = getConfig()

var database = "measures"

// Read the configuration file
func getConfig() (string, string, string, string, string, string, string) {
	file, err := os.Open("/etc/actSoft/conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	a := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a = append(a, scanner.Text())
		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return a[0], a[1], a[2], a[3], a[4], a[5], a[6]
}

// WriteDB Creates batches of points from the actuator responses
func WriteDB(channelResponse chan [38]byte, channelBatches chan client.BatchPoints) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "ms",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"serialNumber": serialNumber, "area": area, "importance": importance, "actFactory": factory, "model": model, "valveTag": tagNumber, "actName": name}
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
				Database:  database,
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
		Addr:               "http://localhost:8087",
		Username:           "bernard",
		Password:           "bernardcontrols",
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
