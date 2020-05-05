package tcp

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ninnemana/drudge"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"github.com/ninnemana/vinyl/pkg/vinyl/firestore"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	tcpAddr = ":8080"
	rpcAddr = ":8081"
)

var (
	rootPattern = runtime.MustPattern(runtime.NewPattern(1, []int{1, 0}, []string(nil), "", runtime.AssumeColonVerbOpt(true)))
	options     = drudge.Options{
		BasePath: "/",
		Addr:     tcpAddr,
		RPC: drudge.Endpoint{
			Network: "tcp",
			Addr:    rpcAddr,
		},
		SwaggerDir: "openapi",
		Handlers: []drudge.Handler{
			// authMiddleware,
			vinyl.RegisterVinylHandler,
		},
		OnRegister: firestore.Register,
	}
)

func Serve() error {
	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}

	// Create and install Jaeger export pipeline
	_, flush, err := jaeger.NewExportPipeline(
		jaeger.WithAgentEndpoint(os.Getenv("JAEGER_AGENT_HOST")+":"+os.Getenv("JAEGER_AGENT_PORT")),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: cfg.ServiceName,
			Tags: []core.KeyValue{
				key.String("exporter", "jaeger"),
			},
		}),
		jaeger.RegisterAsGlobal(),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		return fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	defer flush()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cert, err := tls.LoadX509KeyPair(wd+"/localhost/cert.pem", wd+"/localhost/key.pem")
	if err != nil {
		return fmt.Errorf("failed to load certificate: %w", err)
	}

	options.Certificates = []tls.Certificate{cert}

	options.ServiceName = cfg.ServiceName

	return drudge.Run(context.Background(), options)
}

func authMiddleware(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	fmt.Println("register auth")
	// mux.Handle("POST", rootPattern, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
	// 	ctx, cancel := context.WithCancel(req.Context())
	// 	defer cancel()

	// 	span, _ := opentracing.StartSpanFromContext(ctx, "auth.middleware")
	// 	defer span.Finish()
	// 	fmt.Println("auth")

	// 	// inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
	// 	// rctx, err := runtime.AnnotateContext(ctx, mux, req)
	// 	// if err != nil {
	// 	// 	runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
	// 	// 	return
	// 	// }

	// })

	return nil
}
