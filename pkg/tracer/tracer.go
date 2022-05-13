package tracer

//
//import (
//	"context"
//	"errors"
//	"fmt"
//
//	"contrib.go.opencensus.io/exporter/jaeger"
//	"contrib.go.opencensus.io/exporter/stackdriver"
//	"go.opencensus.io/trace"
//)
//
//type Exporter string
//
//type Config struct {
//	ServiceName string
//	Exporter    Exporter
//	Attributes  map[string]string
//}
//
//var (
//	StackDriver Exporter = "stackdriver"
//	Jaeger      Exporter = "jaeger"
//
//	GCPProjectID            string = "projectID"
//	JaegerAgentEndpoint     string = "jaegerAgent"
//	JaegerCollectorEndpoint string = "jaegerCollector"
//
//	ErrorAttribute string = "error"
//
//	ErrInvalidAttributes      = errors.New("invalid tracing identifier")
//	ErrInvalidProject         = errors.New("invalid project identifier")
//	ErrInvalidJaegerAgent     = errors.New("invalid Jaeger agent endpoint")
//	ErrInvalidJaegerCollector = errors.New("invalid Jaeger collector endpoint")
//
//	unknown int32 = trace.StatusCodeUnknown
//)
//
//func Init(cfg Config) error {
//	var exporter trace.Exporter
//
//	switch cfg.Exporter {
//	case Jaeger:
//		if cfg.Attributes == nil {
//			return ErrInvalidAttributes
//		}
//
//		if cfg.Attributes[JaegerAgentEndpoint] == "" {
//			return ErrInvalidJaegerAgent
//		}
//
//		if cfg.Attributes[JaegerCollectorEndpoint] == "" {
//			return ErrInvalidJaegerCollector
//		}
//
//		// agentEndpointURI := "localhost:6831"
//		// collectorEndpointURI := "http://localhost:14268/api/traces"
//
//		je, err := jaeger.NewExporter(jaeger.Options{
//			AgentEndpoint:     cfg.Attributes[JaegerAgentEndpoint],
//			CollectorEndpoint: cfg.Attributes[JaegerCollectorEndpoint],
//			ServiceName:       cfg.ServiceName,
//		})
//		if err != nil {
//			return fmt.Errorf("failed to create new exporter: %w", err)
//		}
//
//		exporter = je
//	case StackDriver:
//		if cfg.Attributes == nil || cfg.Attributes[GCPProjectID] == "" {
//			return ErrInvalidProject
//		}
//
//		e, err := stackdriver.NewExporter(stackdriver.Options{
//			ProjectID: cfg.Attributes[GCPProjectID],
//		})
//		if err != nil {
//			return fmt.Errorf("failed to create new exporter: %w", err)
//		}
//
//		exporter = e
//	}
//
//	if exporter == nil {
//		return nil
//	}
//
//	trace.RegisterExporter(exporter)
//	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
//
//	return nil
//}
//
//type ErrorConfig struct {
//	Error      error
//	Code       int32
//	Message    string
//	Attributes []trace.Attribute
//}
//
//func RecordError(ctx context.Context, e ErrorConfig) {
//	span := trace.FromContext(ctx)
//	if span == nil {
//		return
//	}
//
//	if e.Error != nil {
//		span.AddAttributes(trace.StringAttribute("error", e.Error.Error()))
//	}
//
//	if e.Code == 0 {
//		e.Code = trace.StatusCodeUnknown
//	}
//
//	span.AddAttributes(e.Attributes...)
//	span.SetStatus(trace.Status{
//		Code:    e.Code,
//		Message: e.Message,
//	})
//}
