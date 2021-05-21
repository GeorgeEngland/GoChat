package main

import (
	"fmt"
	"log"
	"net"

	"github.com/GeorgeEngland/GoChat/chat"
)

func main() {
	s := chat.NewServer()
	go s.Run()
	fmt.Println("hello")
	port := 8888
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("can't start server: %s", err.Error())
	}
	defer listener.Close()
	log.Printf(":Started Server:%v", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error Listening to TCP: %s", err.Error())
			continue
		}
		go s.NewClient(conn)
	}
}
