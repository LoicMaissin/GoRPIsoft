package main

import (
"log"
"encoding/hex"
"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/tty.IBC96342-01001-Bluetoot", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	h, err := hex.DecodeString("011e001f00")
	n, err := s.Write(h)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", buf[:n])
}