package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/tty.IBC96342-01001-Bluetoot", Baud: 115200, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	for count := 0; count < 100; count++ {
		time.Sleep(time.Second)
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

}
