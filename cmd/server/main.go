package main

import (
	"github.com/ninnemana/vinyl/pkg/tcp"
)

func main() {
	if err := tcp.Serve(); err != nil {
		panic(err)
		// log.Fatalf("failed to run vinyl service: %v", err)
	}
}
