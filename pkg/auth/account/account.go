package account

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/codes"

	"github.com/gorilla/mux"
	"github.com/ninnemana/go-discogs"
	"github.com/ninnemana/tracelog"
	"github.com/ninnemana/vinyl/pkg/auth"
	"github.com/ninnemana/vinyl/pkg/users"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	Entity = "auth"
)

var (
	ErrInvalidLogger = errors.New("the provided logger was not valid")

	tracer = otel.Tracer("pkg/auth/account")
)

type Service struct {
	log           *tracelog.TraceLogger
	initTimestamp time.Time
	hostname      string
	users         users.UsersServer
	// tokenizer     auth.Tokenizer
	redirectURL  string
	clientID     string
	clientSecret string
	discogs      discogs.Discogs
	auth.UnimplementedAuthenticationServer
}

type Config struct {
	Logger       *tracelog.TraceLogger
	UserService  users.UsersServer
	Tokenizer    auth.Tokenizer
	Hostname     string
	RedirectURL  string
	ClientID     string
	ClientSecret string
}

func New(ctx context.Context, cfg Config) (*Service, error) {
	if cfg.Logger == nil {
		return nil, ErrInvalidLogger
	}

	return &Service{
		log:           cfg.Logger,
		hostname:      cfg.Hostname,
		initTimestamp: time.Now().UTC(),
		users:         cfg.UserService,
		redirectURL:   cfg.RedirectURL,
		clientID:      cfg.ClientID,
		clientSecret:  cfg.ClientSecret,
	}, nil
}

func (s *Service) Middleware() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}

func (s *Service) Register(rpc *grpc.Server) error {
	s.log.Debug(
		"register RPC service",
		zap.Any("info", rpc.GetServiceInfo()),
	)

	auth.RegisterAuthenticationServer(rpc, s)

	return nil
}

func (s *Service) Route() string {
	return "/" + Entity
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := mux.CurrentRoute(r)
	if route == nil {
		s.log.Error("no route found", zap.String("path", r.URL.Path))
		return
	}

	sub := route.Subrouter()
	sub.HandleFunc("/health", s.HealthHandler)
	sub.HandleFunc("/token", s.TokenHandler)
	sub.HandleFunc("/callback", s.CallbackHandler)

	sub.ServeHTTP(w, r)
}

func (s *Service) Authenticate(ctx context.Context, r *auth.AuthRequest) (*auth.AuthResponse, error) {
	ctx, span := tracer.Start(ctx, "auth/account.Authenticate")
	defer span.End()

	userSvc, ok := s.users.(users.)

	user, err := s.users.Get(ctx, &users.GetParams{
		Email: r.GetEmail(),
	})
	switch err {
	case nil:
	case users.ErrNotFound:
		span.SetStatus(codes.Ok, "user not found")

		return nil, err
	default:
		span.SetStatus(codes.Error, "failed to fetch user")
		span.RecordError(err)

		return nil, fmt.Errorf("failed to fetch user record: %w", err)
	}

	switch s.users.Authenticate(ctx, &users.AuthenticateParams{
		ID:       user.GetId(),
		Password: r.Password,
	}) {
	case nil:
	case users.ErrNotFound:
		span.SetStatus(codes.Ok, "user not found")

		return nil, err
	case users.ErrNotAuthorized:
		span.SetStatus(codes.Ok, "user not authorized")

		return nil, err
	default:
		span.SetStatus(codes.Error, "failed to fetch credentials")
		span.RecordError(err)

		return nil, fmt.Errorf("failed to fetch user credentials: %w", err)
	}

	return &auth.AuthResponse{
		Token: "",
		User:  user,
	}, nil
}

func (s *Service) Callback(ctx context.Context, r *auth.CallbackRequest) (*auth.CallbackResponse, error) {
	return nil, nil
}

func (s *Service) Health(ctx context.Context, r *auth.HealthRequest) (*auth.HealthResponse, error) {
	return &auth.HealthResponse{
		Uptime:  s.initTimestamp.UTC().Format(time.UnixDate),
		Machine: s.hostname,
	}, nil
}

func (s *Service) HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "auth/account.Callback")
	defer span.End()

	resp, err := s.Health(ctx, &auth.HealthRequest{})
	if err != nil {
		http.Error(w, "failed to query health: "+err.Error(), 500)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to marshal health status: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		http.Error(w, "failed to write health status: "+err.Error(), 500)
		return
	}
}

func (s *Service) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "auth/account.Callback")
	defer span.End()

	http.Error(w, "Not supported", http.StatusNotImplemented)
}

func (s *Service) TokenHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "auth/account.TokenHandler")
	defer span.End()

	var req auth.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	s.Authenticate(r.Context(), &req)
}
