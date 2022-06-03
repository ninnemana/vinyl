package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/ninnemana/tracelog"
	httpserver "github.com/ninnemana/vinyl/pkg/http"
	"github.com/ninnemana/vinyl/pkg/router"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Config struct {
	ProjectID               string `env:"GCP_PROJECT_ID"`
	Tracer                  string `env:"TRACER_EXPORTER" envDefault:"jaeger"`
	JaegerAgentEndpoint     string `env:"JAEGER_AGENT"`
	JaegerCollectorEndpoint string `env:"JAEGER_COLLECTOR"`
	LogLevel                string `env:"LOG_LEVEL" envDefault:"info"`
}

const (
	instrumentationName    = "vinyltap-api"
	instrumentationVersion = "v0.1.0"
)

var (
	tracer = otel.GetTracerProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(instrumentationVersion),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
)

func installPipeline(ctx context.Context) (func(context.Context) error, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in an Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(instrumentationName),
		)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("failed to load required arguments: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	logConfig := zap.NewProductionConfig()
	lvl, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err != nil {
		fmt.Printf("failed to parse log level: %v\n", err)
		os.Exit(1)
	}

	logConfig.Level = lvl

	l, err := logConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	flush, err := installPipeline(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer flush(ctx)

	ctx, span := tracer.Start(ctx, "sample")
	defer span.End()

	lg := tracelog.NewLogger(
		tracelog.WithLogger(l),
	)

	//if err := tracer.Init(tracer.Config{
	//	Exporter:    tracer.Exporter(cfg.Tracer),
	//	ServiceName: "vinyltap",
	//	Attributes: map[string]string{
	//		tracer.GCPProjectID:            cfg.ProjectID,
	//		tracer.JaegerAgentEndpoint:     cfg.JaegerAgentEndpoint,
	//		tracer.JaegerCollectorEndpoint: cfg.JaegerCollectorEndpoint,
	//	},
	//}); err != nil {
	//	fmt.Printf("failed to create tracer: %v\n", err)
	//	os.Exit(1)
	//}

	server, err := httpserver.New(lg)
	if err != nil {
		fmt.Printf("failed to run vinyl service: %v\n", err)
		os.Exit(1)
	}

	if err := router.Initialize(lg); err != nil {
		fmt.Printf("failed to run vinyl service: %v\n", err)
		os.Exit(1)
	}

	if err := server.Serve(); err != nil {
		fmt.Printf("failed to run vinyl service: %v\n", err)
		os.Exit(1)
	}
}
