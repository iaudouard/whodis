package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	port = 42069
	host = "localhost"
)

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	defer conn.Close()
	for {
		fmt.Print("Enter command: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write([]byte(scanner.Text()))
		if err != nil {
			log.Fatalf("failed to write to server: %v", err)
		}

		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			log.Fatalf("failed to read from server: %v", err)
		}
		fmt.Println(string(b[:n]))
	}
}
