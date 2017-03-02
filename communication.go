package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
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
	c := &serial.Config{Name: "/dev/tty.IBC96342-01001-Bluetoot", Baud: 115200, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	for count := 0; count < 1; count++ {
		//	time.Sleep(time.Second / 15)
		x := actuatorInfo(getAll(s))
		fmt.Print(x)
	}
}

func getAll(s *serial.Port) []byte {
	// Requête lire toutes les données
	h, _ := hex.DecodeString("011e001F00")
	_, err := s.Write(h)
	if err != nil {
		log.Fatal(err)
	}
	// Reads exactly 38 bytes
	reader := bufio.NewReader(s)
	reply, err := reader.Peek(38)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
	return reply
}

func writeDB(d []byte) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               "https://localhost:8086",
		Username:           os.Getenv("INFLUX_USER"),
		Password:           os.Getenv("INFLUX_PSSWD"),
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Fatal(err)
	} else {
		println("Connected :D")
	}
	defer c.Close() //ensure c is closed after function return

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
	fields := actuatorInfo(d)
	pt, err := client.NewPoint("measures", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
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

func actuatorInfo(response []byte) map[string]interface{} {
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
		"Auxiliary24VFault": nthBit(response[10], 6),
		"TooManyStarts":     nthBit(response[10], 7),

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
