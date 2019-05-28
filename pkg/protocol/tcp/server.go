package tcp

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ninnemana/vinyl/pkg/telemetry"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"github.com/ninnemana/vinyl/pkg/vinyl/datastore"
)

type Server struct {
	log           *zap.Logger
	metricFlush   func()
	context       context.Context
	contextCancel func()
	server        *http.Server
	grpc          *grpc.Server
	vinyls        vinyl.VinylServer
	router        http.Handler
	gwRegFuncs    []gwRegFunc
	gateway       *gwruntime.ServeMux
	address       string
}

// When starting to listen, we will register gateway functions
type gwRegFunc func(ctx context.Context, mux *gwruntime.ServeMux, endpoint string, opts []grpc.DialOption) error

func NewServer(ctx context.Context, projectID, listenAddr string, l *zap.Logger) (*Server, error) {
	metrics, err := telemetry.Register(projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register metrics")
	}

	// GRPC Interceptors
	streamInterceptors := []grpc.StreamServerInterceptor{}
	unaryInterceptors := []grpc.UnaryServerInterceptor{}

	// GRPC Server Options
	serverOptions := []grpc.ServerOption{
		grpc_middleware.WithStreamServerChain(streamInterceptors...),
		grpc_middleware.WithUnaryServerChain(unaryInterceptors...),
	}

	// Create gRPC Server
	grpc := grpc.NewServer(serverOptions...)
	// Register reflection service on gRPC server (so people know what we have)
	reflection.Register(grpc)

	service, err := datastore.New(ctx, l, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create vinyl service")
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/openapi",
		http.StripPrefix(
			"/openapi/",
			http.FileServer(http.Dir("openapi")),
		),
	)

	// Setup the GRPC gateway
	grpcGatewayJSONpbMarshaler := gwruntime.JSONPb(jsonpb.Marshaler{
		EnumsAsInts:  true,
		EmitDefaults: false,
		// OrigName:     ,
	})

	gwmux := gwruntime.NewServeMux(
		gwruntime.WithMarshalerOption(
			gwruntime.MIMEWildcard,
			&grpcGatewayJSONpbMarshaler,
		),
	)

	ctx, cancel := context.WithCancel(ctx)

	s := &Server{
		address:       listenAddr,
		log:           l,
		context:       ctx,
		contextCancel: cancel,
		vinyls:        service,
		gwRegFuncs:    make([]gwRegFunc, 0),
		metricFlush:   metrics,
		grpc:          grpc,
		router:        mux,
		gateway:       gwmux,
	}

	s.server = &http.Server{
		Handler: &Handler{
			server: s,
		},
		ErrorLog: log.New(&errorLogger{l: s.log}, "", 0),
	}

	vinyl.RegisterVinylServer(s.grpc, service)
	s.gwReg(vinyl.RegisterVinylHandlerFromEndpoint)

	return s, nil
}

func (s *Server) ListenAndServe() error {
	defer s.Flush()

	// Listen
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return errors.Wrap(err, "could not establish TCP connection")
	}

	grpcGatewayDialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	s.server.Handler = h2c.NewHandler(s.server.Handler, &http2.Server{})

	// Register all the GRPC gateway functions
	for _, gwrf := range s.gwRegFuncs {
		err = gwrf(context.Background(), s.gateway, listener.Addr().String(), grpcGatewayDialOptions)
		if err != nil {
			return errors.Wrap(err, "could not register HTTP/gRPC gateway")
		}
	}

	s.log.Info("API Listening", zap.String("address", listener.Addr().String()))
	return errors.Wrap(s.server.Serve(listener), "fell out TCP listener")
}

func (s *Server) Flush() {
	if s.contextCancel != nil {
		s.contextCancel()
	}

	if s.metricFlush != nil {
		s.metricFlush()
	}
}

func (s *Server) ServeCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// gwReg will save a gateway registration function for later when the server is started
func (s *Server) gwReg(gwrf gwRegFunc) {
	s.gwRegFuncs = append(s.gwRegFuncs, gwrf)
}

// errorLogger is used for logging errors from the server
type errorLogger struct {
	l *zap.Logger
}

// ErrorLogger implements an error logging function for the server
func (el *errorLogger) Write(b []byte) (int, error) {
	el.l.Error("server failure", zap.Binary("body", b))
	return len(b), nil
}
