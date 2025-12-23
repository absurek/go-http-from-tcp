package main

import (
	"fmt"
	"log"
	"net"

	"github.com/absurek/go-http-from-tcp/internal/constants"
	"github.com/absurek/go-http-from-tcp/internal/request"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("Connection accepted")

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error: request from conn: %v", err)
	}

	fmt.Println("Request line:")
	fmt.Println("- Method:", req.RequestLine.Method)
	fmt.Println("- Target:", req.RequestLine.RequestTarget)
	fmt.Println("- Version:", req.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for key, value := range req.Headers {
		fmt.Printf("- %s: %s\n", key, value)
	}
	fmt.Println()

	log.Println("Connection closed")
}

func main() {
	log.Printf("Setting up tcp listener on %s", constants.Addr)
	listener, err := net.Listen("tcp", constants.Addr)
	if err != nil {
		log.Fatalf("Error: setting up tcp listener: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error: accepting connection: %v", err)
			continue
		}

		handleConnection(conn)
	}
}
