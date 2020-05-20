package main

import (
	"fmt"
	"os"

	httpserver "github.com/ninnemana/vinyl/pkg/http"
	"github.com/ninnemana/vinyl/pkg/log"
	"github.com/ninnemana/vinyl/pkg/router"
)

func main() {
	zlg, err := log.Init()
	if err != nil {
		fmt.Printf("failed to create logger: %v", err)
		os.Exit(1)
	}

	server, err := httpserver.New(zlg)
	if err != nil {
		fmt.Printf("failed to run vinyl service: %v", err)
		os.Exit(1)
	}

	if err := router.Initialize(zlg); err != nil {
		fmt.Printf("failed to run vinyl service: %v", err)
		os.Exit(1)
	}

	if err := server.Serve(); err != nil {
		fmt.Printf("failed to run vinyl service: %v", err)
		os.Exit(1)
	}
}
