package tcp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	// lg "log"
	"mime"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	// "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	// "go.opencensus.io/plugin/ocgrpc"
	"github.com/soheilhy/cmux"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/ninnemana/vinyl/pkg/log"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"github.com/ninnemana/vinyl/pkg/vinyl/datastore"
)

const (
	// rpcAddr  = "vinyltap.alexninneman.com:3000"
	// httpAddr = "vinyltap.alexninneman.com:8080"
	tcpAddr = "vinyltap.alexninneman.com:8000"
)

var (
	certDir = path.Join(os.Getenv("GOPATH"), "src/github.com/ninnemana/vinyl/certs")
)

func getCertificate() (*tls.Certificate, *x509.CertPool, error) {
	certDir := path.Join(os.Getenv("GOPATH"), "src/github.com/ninnemana/vinyl/certs")
	certFile, err := os.Open(path.Join(certDir, "vinyltap.alexninneman.com.crt"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open certificate")
	}

	cert, err := ioutil.ReadAll(certFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read certificate")
	}

	keyFile, err := os.Open(path.Join(certDir, "vinyltap.alexninneman.com.key"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open certificate key")
	}

	key, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read certificate key")
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load certificate")
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		return nil, nil, errors.Wrap(err, "failed to parse certificate")
	}

	return &pair, certPool, nil
}

func getCredentials() (credentials.TransportCredentials, error) {
	// certFile, err := os.Open(path.Join(certDir, "vinyltap.alexninneman.com.crt"))
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to open certificate")
	// }

	// keyFile, err := os.Open(path.Join(certDir, "server.key"))
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to open certificate key")
	// }

	creds, err := credentials.NewServerTLSFromFile(
		path.Join(certDir, "vinyltap.alexninneman.com.crt"),
		path.Join(certDir, "vinyltap.alexninneman.com.key"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transport credentials")
	}

	return creds, nil
}

func getClientCredentials() (credentials.TransportCredentials, error) {
	certDir := path.Join(os.Getenv("GOPATH"), "src/github.com/ninnemana/vinyl/certs")
	// certFile, err := os.Open(path.Join(certDir, "vinyltap.alexninneman.com.crt"))
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to open certificate")
	// }

	creds, err := credentials.NewClientTLSFromFile(
		path.Join(certDir, "vinyltap.alexninneman.com.crt"),
		"",
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client certificate")
	}

	return creds, nil
}

func telemetry() (func(), error) {
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "vinyl-registry",
		// MetricPrefix helps uniquely identify your metrics.
		MetricPrefix: "vinyl-registry",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stats exporter")
	}

	// Register it as a metrics exporter
	view.RegisterExporter(sd)
	view.SetReportingPeriod(60 * time.Second)

	// Register it as a trace exporter
	trace.RegisterExporter(sd)

	return sd.Flush, nil
}

func Serve() error {
	// _, pair, err := getCertificate()
	// if err != nil {
	// 	return err
	// }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := log.Init(); err != nil {
		return errors.Wrap(err, "Failed to create logger")
	}

	flush, err := telemetry()
	if err != nil {
		return err
	}
	defer flush()

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		return errors.Wrap(err, "failed to start TCP listener")
	}

	m := cmux.New(lis)
	// Using MatchWithWriters/SendSettings is a major performance hit (around 15%).
	// Per the cmux documentation, you have to do this for grpc-java.
	// If only using golang, you don't need this, but probably not
	// great to assume what the calling languages are.
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	server, err := datastore.New(context.Background(), log.Logger, "vinyl-registry")
	if err != nil {
		return errors.Wrap(err, "failed to create vinyl service")
	}

	rpc := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			log.UnaryServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			log.StreamServerInterceptor(),
		),
	)
	vinyl.RegisterVinylServer(rpc, server)
	go func() {
		if err := rpc.Serve(httpL); err != nil {
			log.Logger.Fatal("failed to start HTTP server listening", zap.Error(err))
		}
	}()

	grpcS := grpc.NewServer()

	// pb.RegisterGreeterServer(grpcS, srv)
	// Register reflection service on gRPC server.
	// reflection.Register(grpcS)
	go func() {
		if err := grpcS.Serve(grpcL); err != nil {
			log.Logger.Fatal("failed to start RPC server", zap.Error(err))
		}
	}()

	return errors.Wrap(m.Serve(), "fell out of service traffic")
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Content-Type"))
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			fmt.Println(r.ProtoMajor, "grpc handler")
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func serveSwagger(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	mux.Handle(
		"/openapi/",
		http.StripPrefix(
			"/openapi/",
			http.FileServer(http.Dir("openapi")),
		),
	)
}
