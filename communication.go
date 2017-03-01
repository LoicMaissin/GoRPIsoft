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
	c := &serial.Config{Name: "/dev/tty.IBC96342-01001-Bluetoot", Baud: 115200, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	for count := 0; count < 100; count++ {
		//	time.Sleep(time.Second / 15)
		getAll(s)
	}
}

func getAll(s *serial.Port) {
	// Requête lire toutes les données
	h, err := hex.DecodeString("011e001F00")
	_, err = s.Write(h)
	if err != nil {
		log.Fatal(err)
	}
	// Reads exactly 17 bytes
	reader := bufio.NewReader(s)
	reply, err := reader.Peek(32)
	if err != nil {
		panic(err)
	}
	log.Println("Ouverture de la vanne (%)", reply[9])
	log.Println("Couple subi (%)", reply[8])
	writeDB(reply)
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
	fields := map[string]interface{}{
		"currentTorque": d[8],
		"Opening":       d[9],
	}

	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}
