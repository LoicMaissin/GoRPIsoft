package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/tarm/serial"
)

const (
	myDB = "test"
)

func main() {

	channelResponse := make(chan [38]byte)
	channelBatches := make(chan client.BatchPoints)

	// Launch the Database client
	go writeDB(channelResponse, channelBatches)
	go sendDB(channelBatches)

	c := &serial.Config{
		Name:        "/dev/tty.IBC96342-01001-Bluetoot",
		Baud:        115200,
		ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal("Opening")
	}
	for count := 0; count < 500; count++ {
		channelResponse <- getAll(s)
	}
	close(channelResponse)
	time.Sleep(time.Second * 3)
}

func getAll(s *serial.Port) [38]byte {
	// Requête lire toutes les données
	h, _ := hex.DecodeString("011e001F00")
	_, err := s.Write(h)
	for err != nil {
		_, err = s.Write(h)
		time.Sleep(time.Second)
	}
	// Reads exactly 38 bytes
	reader := bufio.NewReader(s)
	reply, err := reader.Peek(38)
	if err != nil {
		log.Println("Error reading buffer")
		log.Fatal(err)
	}

	var res [38]byte
	copy(res[:], reply)
	return res
}

func writeDB(channelResponse chan [38]byte, channelBatches chan client.BatchPoints) {

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  myDB,
		Precision: "ms",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"actuator": "0"}
	last := [38]byte{}
	lastTime := time.Now()
	count := 0
	for response := range channelResponse {
		if last == response {
			lastTime = time.Now()
		} else {
			fields := actuatorInfo(last)
			pt, errPt := client.NewPoint("measures", tags, fields, lastTime)
			if errPt != nil {
				log.Fatal(errPt)
			}
			bp.AddPoint(pt)

			fields = actuatorInfo(response)
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
	fields := actuatorInfo(last)
	pt, errPt := client.NewPoint("measures", tags, fields, lastTime)
	if errPt != nil {
		log.Fatal(errPt)
	}
	bp.AddPoint(pt)
	channelBatches <- bp
	time.Sleep(time.Second * 2)

	close(channelBatches)
}

func sendDB(channelBatches chan client.BatchPoints) {
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

func nthBit(x byte, n uint) bool {
	if x>>n&1 == byte(0) {
		return false
	}
	return true
}

func addBytes(v []byte) int {
	// The number equal to v concatenated, with v[0] MSB
	res := 0
	for _, x := range v {
		res = res<<8 + int(x)
	}
	return res
}

func actuatorInfo(response [38]byte) map[string]interface{} {
	fields := map[string]interface{}{
		"isOpened":                          nthBit(response[5], 0),
		"isClosed":                          nthBit(response[5], 1),
		"torqueLimiterActionOpenDirection":  nthBit(response[5], 2),
		"torqueLimiterActionCloseDierction": nthBit(response[5], 3),
		"selectorToLocalPosition":           nthBit(response[5], 4),
		"selectorToRemotePosition":          nthBit(response[5], 5),
		"selectorToOffPosition":             nthBit(response[5], 6),
		"powerOn":                           nthBit(response[5], 7),

		"actOpening":               nthBit(response[6], 0),
		"actClosing":               nthBit(response[6], 1),
		"handwheelAction":          nthBit(response[6], 2),
		"ESDCommand":               nthBit(response[6], 3),
		"actRunning":               nthBit(response[6], 4),
		"actFault":                 nthBit(response[6], 5),
		"positionSensorPowerFault": nthBit(response[6], 6),
		"torqueSensorPowerFault":   nthBit(response[6], 7),

		"lockedMotorOpen":      nthBit(response[7], 0),
		"lockedMotorClose":     nthBit(response[7], 1),
		"motorThermalOverload": nthBit(response[7], 2),
		"lostPhase":            nthBit(response[7], 3),
		"overtravelAlarm":      nthBit(response[7], 4),
		"directionOpenAlarm":   nthBit(response[7], 5),
		"directionCloseAlarm":  nthBit(response[7], 6),
		"batteryLow":           nthBit(response[7], 7),

		"runningTorque": response[8],
		"actPosition":   response[9],

		"indication1":       nthBit(response[10], 0),
		"indication2":       nthBit(response[10], 1),
		"indication3":       nthBit(response[10], 2),
		"indication4":       nthBit(response[10], 3),
		"indication5":       nthBit(response[10], 4),
		"valveJammed":       nthBit(response[10], 5),
		"auxiliary24VFault": nthBit(response[10], 6),
		"tooManyStarts":     nthBit(response[10], 7),

		"pumping":              nthBit(response[11], 0),
		"confMemFault":         nthBit(response[11], 1),
		"activityMemFault":     nthBit(response[11], 2),
		"baseMemFault":         nthBit(response[11], 3),
		"stopMidTravel":        nthBit(response[11], 4),
		"lostSignal":           nthBit(response[11], 5),
		"partialStrokeRunning": nthBit(response[11], 6),
		"partialStrokeFault":   nthBit(response[11], 7),

		"openBreakoutMaxTorque": response[12],
		"closeTightMaxTorque":   response[13],
		"openingMaxTorque":      response[14],
		"closingMaxTorque":      response[15],
		"startsLast12h":         addBytes(response[16:18]),
		"totalStarts":           addBytes(response[18:22]),
		"totalRunningTime":      addBytes(response[22:26]),
		"partialStarts":         addBytes(response[26:30]),
		"partialRunningTime":    addBytes(response[30:34]),
		"actPosition(per mil)":  addBytes(response[34:36]),
	}
	return fields
}
