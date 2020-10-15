package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	httpserver "github.com/ninnemana/vinyl/pkg/http"
	"github.com/ninnemana/vinyl/pkg/log"
	"github.com/ninnemana/vinyl/pkg/router"
	"github.com/ninnemana/vinyl/pkg/tracer"
)

var (
	projectID = os.Getenv("GCP_PROJECT_ID")
)

type Config struct {
	ProjectID               string `env:"GCP_PROJECT_ID"`
	Tracer                  string `env:"TRACER_EXPORTER" envDefault:"jaeger"`
	JaegerAgentEndpoint     string `env:"JAEGER_AGENT"`
	JaegerCollectorEndpoint string `env:"JAEGER_COLLECTOR"`
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("failed to load required arguments: %v\n", err)
		os.Exit(1)
	}

	zlg, err := log.Init()
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}

	if err := tracer.Init(tracer.Config{
		Exporter:    tracer.Exporter(cfg.Tracer),
		ServiceName: "vinyltap",
		Attributes: map[string]string{
			tracer.GCPProjectID:            cfg.ProjectID,
			tracer.JaegerAgentEndpoint:     cfg.JaegerAgentEndpoint,
			tracer.JaegerCollectorEndpoint: cfg.JaegerCollectorEndpoint,
		},
	}); err != nil {
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
