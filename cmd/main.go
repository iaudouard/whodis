package main

import (
	"log"

	"github.com/iaudouard/whodis/tcp"
)

func main() {
	s, err := tcp.NewTCPServer(42069)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	s.Listen()
}
