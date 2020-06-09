package httpserver

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	server = &Server{
		log:  zap.NewNop(),
		rpc:  grpc.NewServer(),
		http: &http.Server{Handler: mux.NewRouter().StrictSlash(true)},
	}
)

type Registerer func(*grpc.Server, interface{})

type Handler interface {
	http.Handler
	Route() string
	Register(*grpc.Server) error
	Middleware() []mux.MiddlewareFunc
}

type Server struct {
	log      *zap.Logger
	http     *http.Server
	rpc      *grpc.Server
	listener net.Listener
}

func New(log *zap.Logger) (*Server, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}

	server.log = log
	server.listener = listener

	return server, nil
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
	sub := mx.PathPrefix(h.Route()).Subrouter()
	sub.Use(h.Middleware()...)
	sub.NewRoute().Handler(h)

	return nil
}

func RegisterRPC(r Registerer, svc interface{}) error {
	r(server.rpc, svc)

	return nil
}

func (s *Server) Serve() error {
	s.log.Debug("starting http/rpc listeners")
	m := cmux.New(server.listener)
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	mx, ok := s.http.Handler.(*mux.Router)
	if !ok {
		return fmt.Errorf("router was not expected type: %T", server.http.Handler)
	}

	spa := spaHandler{staticPath: "./ui/dist", indexPath: ""}
	mx.Use(s.Logger)
	mx.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			s.log.Info("options")
		}
		spa.ServeHTTP(w, r)
	})

	s.http.Handler = cors.AllowAll().Handler(mx)

	g := new(errgroup.Group)
	g.Go(func() error { return s.rpc.Serve(grpcListener) })
	g.Go(func() error { return s.http.Serve(httpListener) })
	g.Go(func() error { return m.Serve() })

	return g.Wait()
}

func (s *Server) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("handling route", zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, h.staticPath)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
