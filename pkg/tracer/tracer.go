package tracer

import (
	"fmt"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Init(projectID string) error {
	exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	global.SetTracerProvider(
		sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter)),
	)

	return nil
}
