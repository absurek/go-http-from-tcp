package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/absurek/go-http-from-tcp/internal/constants"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", constants.LocalAddress)
	if err != nil {
		log.Fatalf("Error: resolve udp address at %s: %v", constants.LocalAddress, err)
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatalf("Error: dial udp at %s: %v", address, err)
	}
	defer conn.Close()

	console := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := console.ReadString('\n')
		if err != nil {
			log.Printf("Error: reading input: %v", err)
			continue
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Printf("Error: writing to connection: %v", err)
		}
	}
}
