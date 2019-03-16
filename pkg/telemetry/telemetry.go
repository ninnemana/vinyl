package telemetry

import (
	"fmt"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/pkg/errors"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func Register(projectID string) (func(), error) {
	fmt.Println(projectID)
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
		// MetricPrefix helps uniquely identify your metrics.
		MetricPrefix: projectID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stats exporter")
	}

	// Register it as a metrics exporter
	view.RegisterExporter(sd)
	view.SetReportingPeriod(60 * time.Second)

	// Register it as a trace exporter
	trace.RegisterExporter(sd)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return sd.Flush, nil
}
