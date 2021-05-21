package main

import (
	"log"
	"net"

	"github.com/GeorgeEngland/GoChat/chat"
)

func main() {
	s := chat.NewServer()
	go s.Run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("can't start server: %s", err.Error())
	}
	defer listener.Close()
	log.Printf("Started Server: 8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error Listening to TCP: %s", err.Error())
			continue
		}
		go s.NewClient(conn)
	}
}
