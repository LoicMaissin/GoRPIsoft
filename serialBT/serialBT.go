package serialBT

/*
	Handle the Bluetooth connection with the actuator
	Needs *SERIAL to be defined to the location of the serial connection
*/

import (
	"bufio"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

func initSerial() *serial.Port {
	// Open the serial connection
	c := &serial.Config{
		Name:        os.Getenv("SERIAL"),
		Baud:        115200,
		ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	for err != nil {
		log.Println("Opening")
		time.Sleep(time.Second * 30)
		s, err = serial.OpenPort(c)
	}
	log.Println("Connected")
	return s
}

var port = initSerial()

// GetAll requests data from the actuator and return its response
func GetAll() [38]byte {
	// Requests for all the data
	h, _ := hex.DecodeString("011e001F00")
	_, errWrite := port.Write(h)
	if errWrite != nil {
		log.Println(errWrite)
	}
	// Reads exactly 38 bytes
	reader := bufio.NewReader(port)
	reply, err := reader.Peek(37)
	for err != nil {
		log.Println("Error reading buffer")
		port.Close()
		time.Sleep(time.Second * 30)
		port = initSerial()
		reader = bufio.NewReader(port)
		_, errWrite = port.Write(h)
		if errWrite != nil {
			log.Println(errWrite)
		}
		reply, err = reader.Peek(38)
	}
	var res [38]byte
	copy(res[:], reply)
	return res
}
