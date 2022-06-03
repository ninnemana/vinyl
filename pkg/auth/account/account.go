package account

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/codes"

	"github.com/gorilla/mux"
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
	//tokenizer     auth.Tokenizer
	//redirectURL  string
	//clientID     string
	//clientSecret string
	//discogs      discogs.Discogs
	auth.UnimplementedAuthenticationServer
}

type Config struct {
	Logger      *tracelog.TraceLogger
	UserService users.UsersServer
	Hostname    string
	//Tokenizer   auth.Tokenizer
	//RedirectURL  string
	//ClientID     string
	//ClientSecret string
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
		//redirectURL:   cfg.RedirectURL,
		//clientID:      cfg.ClientID,
		//clientSecret:  cfg.ClientSecret,
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

	if r.GetEmail() == "" {
		span.SetStatus(codes.Ok, users.ErrInvalidEmail.Error())

		return nil, users.ErrNotFound
	}

	s.log.Debug("attempting to authenticate", zap.String("email", r.GetEmail()))

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

	resp, err := s.users.Authenticate(ctx, &users.AuthenticateRequest{
		UserID:   user.GetId(),
		Password: r.Password,
	})
	switch err {
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
		Token: resp.GetToken(),
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
		span.SetStatus(codes.Ok, "request body was not valid")
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	s.log.Info("we out here")

	s.log.Debug("attempting to authenticate against service", zap.String("email", req.Email))
	resp, err := s.Authenticate(r.Context(), &req)
	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			span.SetStatus(codes.Error, "failed to encode auth response")
			span.RecordError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case users.ErrNotAuthorized:
		span.SetStatus(codes.Ok, "user is not authorized")
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case users.ErrNotFound:
		span.SetStatus(codes.Ok, "user does not exist")
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		span.SetStatus(codes.Error, "failed to execute auth request")
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
