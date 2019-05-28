package main

import (
	"log"

	"github.com/ninnemana/vinyl/pkg/protocol/tcp"
)

func main() {
	if err := tcp.Serve(); err != nil {
		log.Fatalf("failed to run vinyl service: %v", err)
	}
}