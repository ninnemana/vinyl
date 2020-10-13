package main

import (
	"fmt"
	"os"

	httpserver "github.com/ninnemana/vinyl/pkg/http"
	"github.com/ninnemana/vinyl/pkg/log"
	"github.com/ninnemana/vinyl/pkg/router"
	"github.com/ninnemana/vinyl/pkg/tracer"
)

var (
	projectID = os.Getenv("GCP_PROJECT_ID")
)

func main() {
	fmt.Printf("%+v\n", os.Environ())
	zlg, closer, err := log.Init()
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer closer()

	if err := tracer.Init(projectID); err != nil {
		fmt.Printf("failed to create tracer: %v\n", err)
		os.Exit(1)
	}

	server, err := httpserver.New(zlg)
	if err != nil {
		fmt.Printf("failed to run vinyl service: %v\n", err)
		os.Exit(1)
	}

	if err := router.Initialize(zlg); err != nil {
		fmt.Printf("failed to run vinyl service: %v\n", err)
		os.Exit(1)
	}

	if err := server.Serve(); err != nil {
		fmt.Printf("failed to run vinyl service: %v\n", err)
		os.Exit(1)
	}
}
