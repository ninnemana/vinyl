package main

import (
	"context"
	"log"

	rest "github.com/ninnemana/vinyl/pkg/protocol/http"
)

func main() {
	if err := rest.Start(context.Background(), "localhost:8080", ":8000"); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
