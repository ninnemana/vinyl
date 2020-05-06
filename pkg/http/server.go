package httpserver

import (
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	server = &Server{
		log:  zap.NewNop(),
		rpc:  grpc.NewServer(),
		http: &http.Server{Handler: mux.NewRouter()},
	}
)

type Registerer func(*grpc.Server, interface{})

type Handler interface {
	http.Handler
	Route() string
	Register(*grpc.Server) error
}

type Server struct {
	log      *zap.Logger
	http     *http.Server
	rpc      *grpc.Server
	listener net.Listener
}

func New(log *zap.Logger) error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	server.log = log
	server.listener = listener

	return nil
}

func RegisterHandler(h Handler) error {
	if err := h.Register(server.rpc); err != nil {
		return err
	}

	mx, ok := server.http.Handler.(*mux.Router)
	if !ok {
		return fmt.Errorf(
			"the provided HTTP mux was not expected '%T' type, got '%T'",
			&mux.Router{},
			server.http.Handler,
		)
	}

	server.log.Debug("registering handler", zap.String("route", h.Route()))
	mx.PathPrefix(h.Route()).Handler(h)

	return nil
}

func RegisterRPC(r Registerer, svc interface{}) error {
	r(server.rpc, svc)

	return nil
}

func Serve(log *zap.Logger) error {
	if err := New(log); err != nil {
		return err
	}

	log.Debug("starting http/rpc listeners")
	m := cmux.New(server.listener)
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	g := new(errgroup.Group)
	g.Go(func() error { return server.rpc.Serve(grpcListener) })
	g.Go(func() error { return server.http.Serve(httpListener) })
	g.Go(func() error { return m.Serve() })

	return g.Wait()
}
