package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/absurek/go-http-from-tcp/internal/constants"
)

func readLinesChannel(f io.ReadCloser) <-chan string {
	lineCh := make(chan string)

	go func() {
		defer f.Close()
		defer close(lineCh)

		var line string
		buffer := make([]byte, 8)
		for {
			n, err := f.Read(buffer)
			if err != nil {
				switch {
				case errors.Is(err, io.EOF):
					lineCh <- line + string(buffer[:n])
				default:
					log.Fatalf("Error: %v", err)
				}

				return
			}

			parts := strings.Split(string(buffer[:n]), "\n")
			for i, part := range parts {
				if i < len(parts)-1 {
					lineCh <- line + part
					line = ""
				} else {
					line += part
				}
			}
		}
	}()

	return lineCh
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connection accepted")

	lineCh := readLinesChannel(conn)
	for line := range lineCh {
		fmt.Println(line)
	}

	fmt.Println("Connection closed")
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
