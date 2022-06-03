package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ninnemana/vinyl/pkg/auth"
	"github.com/ninnemana/vinyl/pkg/users"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ninnemana/tracelog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	tracer = otel.Tracer("pkg/users/firestore")
)

const (
	Entity     = "users"
	AuthEntity = "credentials"
)

type Service struct {
	log       *tracelog.TraceLogger
	firestore *firestore.Client
	projectID string
	tokenizer auth.Tokenizer
	users.UnimplementedUsersServer
}

type Authentication struct {
	Password string
	Created  time.Time
	Updated  time.Time
}

func New(ctx context.Context, l *tracelog.TraceLogger, projectID string, tokenizer auth.Tokenizer, opts ...option.ClientOption) (*Service, error) {
	ctx, span := tracer.Start(ctx, "users/firestore.New")
	defer span.End()

	client, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		//tracer.RecordError(ctx, tracer.ErrorConfig{
		//	Error:   err,
		//	Message: "failed to create Firestore client",
		//	Code:    trace.StatusCodeFailedPrecondition,
		//})

		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return &Service{
		l,
		client,
		projectID,
		tokenizer,
		users.UnimplementedUsersServer{},
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

	users.RegisterUsersServer(rpc, s)

	return nil
}

func (s *Service) Route() string {
	return "/" + Entity
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "users/firestore.ServeHTTP")
	defer span.End()

	r = r.WithContext(ctx)

	route := mux.CurrentRoute(r)
	if route == nil {
		span.SetAttributes(attribute.String("path", r.URL.Path))
		span.SetStatus(codes.Error, "no route found")

		s.log.Error("no route found", zap.String("path", r.URL.Path))
		return
	}

	sub := route.Subrouter()
	sub.HandleFunc("/health", s.HealthHandler)
	sub.Methods(http.MethodPost).Path("").HandlerFunc(s.SaveHandler)
	sub.HandleFunc("", s.GetHandler)

	sub.ServeHTTP(w, r)
}

func (s *Service) SaveHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "users/firestore.SaveHandler")
	defer span.End()

	var u users.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		span.SetStatus(codes.Error, "failed to decode request body")
		span.RecordError(err)

		s.log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("email", u.GetEmail()))

	result, err := s.Save(ctx, &u)
	switch err.(type) {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			span.SetStatus(codes.Error, "failed to encode user")
			span.RecordError(err)

			s.log.Error("failed to encode user", zap.Error(err))

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case users.ValidationError:
		s.log.Debug("user failed validation", zap.Error(err))
		// span.SetStatus(codes.Error, "user failed validation")

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		span.SetStatus(codes.Error, "failed to save user")
		span.RecordError(err)

		s.log.Error("failed to save user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "users/firestore.GetHandler")
	defer span.End()

	var params users.GetParams
	switch r.Method {
	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			span.SetStatus(codes.Error, "failed to decode request body")
			span.RecordError(err)
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
	case http.MethodGet:
		var err error
		params, err = users.GetFromQueryString(r.URL.Query())
		if err != nil {
			span.SetStatus(codes.Error, "failed to parse query string parameters")
			span.RecordError(err)
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
	default:
		span.SetStatus(codes.Ok, fmt.Sprintf("provided '%s' HTTP method is unsupported", r.Method))
		http.Error(w, fmt.Sprintf("provided '%s' HTTP method is unsupported", r.Method), http.StatusMethodNotAllowed)

		return
	}

	resp, err := s.Get(ctx, &params)
	if err != nil {
		span.SetStatus(codes.Error, "failed to save user")
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		span.SetStatus(codes.Error, "failed to encode request body")
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (s *Service) HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "users/firestore.HealthHandler")
	defer span.End()

	resp, err := s.Health(ctx, &users.HealthRequest{})
	if err != nil {
		span.SetStatus(codes.Error, "failed to execute health check")
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		span.SetStatus(codes.Error, "failed to encode response")
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (s *Service) Authenticate(ctx context.Context, p *users.AuthenticateRequest) (*users.AuthenticateResponse, error) {
	ctx, span := tracer.Start(ctx, "users/firestore.Authenticate")
	defer span.End()

	s.log.Debug(
		"authenticate user",
		zap.String("userID", p.UserID),
		zap.String("password", p.Password),
	)

	user, err := s.Get(ctx, &users.GetParams{
		Id: p.GetUserID(),
	})
	if err != nil {
		return nil, err
	}

	doc, err := s.firestore.Collection(AuthEntity).Doc(p.UserID).Get(ctx)
	switch status.Code(err) {
	case grpccodes.OK:
		var creds Authentication
		if err := doc.DataTo(&creds); err != nil {
			return nil, fmt.Errorf("failed to decode credential entry: %w", err)
		}

		if !users.ComparePasswords([]byte(creds.Password), []byte(p.Password)) {
			return nil, users.ErrNotAuthorized
		}

		token, err := s.tokenizer.GenerateToken(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to generate auth token: %w", err)
		}

		return &users.AuthenticateResponse{
			Token: token,
		}, nil
	case grpccodes.NotFound:
		span.SetStatus(codes.Ok, "credentials not found for user")
		span.RecordError(err)

		return nil, users.ErrNotFound
	default:
		span.SetStatus(codes.Error, "failed to fetch credentials for user")
		span.RecordError(err)

		return nil, fmt.Errorf("failed to fetch credentials: %w", err)
	}
}

func (s *Service) Health(ctx context.Context, p *users.HealthRequest) (*users.HealthResponse, error) {
	_, span := tracer.Start(ctx, "users/firestore.Health")
	defer span.End()

	return &users.HealthResponse{
		Uptime:  "",
		Machine: "",
	}, nil
}

func (s *Service) Get(ctx context.Context, p *users.GetParams) (*users.User, error) {
	ctx, span := tracer.Start(ctx, "users/firestore.Get")
	defer span.End()

	s.log.Debug(
		"fetch user",
		zap.String("params", p.String()),
	)

	span.SetAttributes(attribute.String("params", p.String()))
	span.AddEvent("fetching user")

	iter := p.Where(s.firestore.Collection(Entity).Query).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			span.SetStatus(codes.Error, "failed to fetch users")
			span.RecordError(err)

			return nil, err
		}

		//data := doc.Data()
		//var user users.User
		//if err := doc.DataTo(&user); err != nil {
		//	span.SetStatus(codes.Error, "failed to parse user")
		//	span.RecordError(err)
		//
		//	return nil, err
		//}

		data := doc.Data()
		user := &users.User{
			Id:      toString(data["id"]),
			Name:    toString(data["name"]),
			Email:   toString(data["email"]),
			Created: toTimestamp(data["created"]),
			Updated: toTimestamp(data["updated"]),
			Deleted: toTimestamp(data["deleted"]),
		}

		if user.Deleted != nil {
			continue
		}

		return user, nil
	}

	span.SetStatus(codes.Ok, "user not found")

	return nil, users.ErrNotFound
}

func (s *Service) Save(ctx context.Context, u *users.User) (*users.User, error) {
	ctx, span := tracer.Start(ctx, "users/firestore.Save")
	defer span.End()

	// note: we're doing the validation here so we
	// can control the response type.
	if err := u.Validate(); err != nil {
		span.SetStatus(codes.Error, "failed to validate user")
		span.RecordError(err)

		return nil, err
	}

	span.SetAttributes(attribute.String("email", u.GetEmail()))

	_, err := s.Get(ctx, &users.GetParams{
		Email: u.GetEmail(),
	})
	switch err {
	case nil:
		span.SetStatus(codes.Ok, "user already exists")

		return nil, users.ErrUserExists
	case users.ErrNotFound:
	default:
		span.SetStatus(codes.Error, "failed to fecth user")
		span.RecordError(err)

		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	// new user workflow
	newUser := u.Id == ""
	if u.Id == "" {
		u.Id = uuid.New().String()
		u.Created = &timestamp.Timestamp{
			Seconds: time.Now().UTC().Unix(),
		}

		pwd, err := users.HashAndSalt([]byte(u.Password))
		if err != nil {
			span.SetStatus(codes.Error, "failed to encode password")
			span.RecordError(err)

			return nil, err
		}

		u.Password = pwd
	}

	defer func() {
		u.Password = ""
		u.AuthenticatedAccounts = nil
	}()

	u.Updated = &timestamp.Timestamp{
		Seconds: time.Now().UTC().Unix(),
	}

	s.log.Debug("saving users", zap.String("id", u.GetId()))

	js, err := json.Marshal(u)
	if err != nil {
		span.SetStatus(codes.Error, "failed to encode user")
		span.RecordError(err)

		s.log.Error("failed to marshal user", zap.Error(err))
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(js, &data); err != nil {
		span.SetStatus(codes.Error, "failed to parse encoding")
		span.RecordError(err)

		return nil, err
	}

	delete(data, "password")

	coll := s.firestore.Collection(Entity)
	if _, err := coll.Doc(u.GetId()).Set(ctx, data); err != nil {
		s.log.Error("failed to save user to store", zap.Error(err))
		span.SetStatus(codes.Error, "failed to save user")
		span.RecordError(err)

		return nil, fmt.Errorf("failed to set user into store: %w", err)
	}

	s.log.Debug("successfully saved user", zap.String("id", u.GetId()))

	if newUser {
		if _, err := s.firestore.Collection(AuthEntity).Doc(u.GetId()).Set(ctx, Authentication{
			Password: u.Password,
			Created:  time.Now().UTC(),
			Updated:  time.Now().UTC(),
		}); err != nil {
			s.log.Error("failed to save user credentials to store", zap.Error(err))
			span.SetStatus(codes.Error, "failed to save user credentials")
			span.RecordError(err)

			return nil, fmt.Errorf("failed to set user into store: %w", err)
		}
	}
	return u, nil
}

func toString(i interface{}) string {
	v, ok := i.(string)
	if !ok {
		return ""
	}

	return v
}

func toTimestamp(i interface{}) *timestamp.Timestamp {
	v, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	if v["seconds"] == nil {
		return nil
	}

	sec, ok := v["seconds"].(float64)
	if !ok {
		return nil
	}

	return &timestamp.Timestamp{
		Seconds: int64(sec),
	}
}
