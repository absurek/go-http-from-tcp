package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func readLinesChannel(f io.ReadCloser) <-chan string {
	lineCh := make(chan string)

	go func() {
		var line string
		buffer := make([]byte, 8)
		for {
			n, err := f.Read(buffer)
			if err != nil {
				switch {
				case errors.Is(err, io.EOF):
					lineCh <- line + string(buffer[:n])
				default:
					log.Fatalf("Error: unexpected error while reading messages.txt: %v", err)
				}

				close(lineCh)
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

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Error: Could not open messages.txt: %v", err)
	}
	defer f.Close()

	lineCh := readLinesChannel(f)
	for line := range lineCh {
		fmt.Println("read:", line)
	}
}
