package grpc

import (
	"context"
	lg "log"
	"net"
	"time"

	"github.com/ninnemana/vinyl/pkg/log"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"github.com/ninnemana/vinyl/pkg/vinyl/datastore"

	"contrib.go.opencensus.io/exporter/stackdriver"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Serve() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := log.Init(); err != nil {
		return errors.Wrap(err, "failed to create logger")
	}

	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "vinyl-registry",
		// MetricPrefix helps uniquely identify your metrics.
		MetricPrefix: "vinyl-registry",
	})
	if err != nil {
		return errors.Wrap(err, "failed to create stats exporter")
	}
	// It is imperative to invoke flush before your main function exits
	defer sd.Flush()

	// Register it as a metrics exporter
	view.RegisterExporter(sd)
	view.SetReportingPeriod(60 * time.Second)

	// Register it as a trace exporter
	trace.RegisterExporter(sd)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return errors.Wrap(err, "failed to create tcp listener")
	}

	defer func() {
		if err := l.Close(); err != nil {
			lg.Fatalf("Failed to close %s %s: %v", "tcp", ":8080", err)
		}
	}()

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			log.UnaryServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			log.StreamServerInterceptor(),
		),
	)

	server, err := datastore.New(context.Background(), log.Logger, "vinyl-registry")
	if err != nil {
		return errors.Wrap(err, "failed to create vinyl service")
	}

	vinyl.RegisterVinylServer(s, server)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	log.Logger.Debug("Serving RPC", zap.String("port", ":8080"))
	return errors.Wrap(s.Serve(l), "fell out of serving traffic")
}
