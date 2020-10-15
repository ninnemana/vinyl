package tracer

import (
	"context"
	"fmt"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

var (
	ErrorAttribute string = "error"

	unknown int32 = trace.StatusCodeUnknown
)

func Init(projectID string) error {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		return fmt.Errorf("failed to create new exporter: %w", err)
	}

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return nil
}

type ErrorConfig struct {
	Error      error
	Code       int32
	Message    string
	Attributes []trace.Attribute
}

func RecordError(ctx context.Context, e ErrorConfig) {
	span := trace.FromContext(ctx)
	if span == nil {
		return
	}

	if e.Error != nil {
		span.AddAttributes(trace.StringAttribute("error", e.Error.Error()))
	}

	if e.Code == 0 {
		e.Code = trace.StatusCodeUnknown
	}

	span.AddAttributes(e.Attributes...)
	span.SetStatus(trace.Status{
		Code:    e.Code,
		Message: e.Message,
	})
}
