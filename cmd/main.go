package main

import (
	"log"

	"github.com/iaudouard/whodis/tcp"
)

const (
	port = 42069
)

func main() {
	s, err := tcp.NewTCPServer(port)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	s.Listen()
}
