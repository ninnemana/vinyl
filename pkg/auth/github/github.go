package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/ninnemana/vinyl/pkg/auth"
	"github.com/ninnemana/vinyl/pkg/users"

	"github.com/google/go-github/v31/github"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

const (
	Entity = "auth"
)

var (
	ErrInvalidLogger = errors.New("the provided logger was not valid")
	httpClient       = &http.Client{
		Timeout: time.Second * 5,
	}
)

type Service struct {
	log           *zap.Logger
	initTimestamp time.Time
	hostname      string
	http          *http.Client
	github        *github.Client
	users         users.UsersServer
	tokenizer     auth.Tokenizer
	redirectURL   string
	clientID      string
	clientSecret  string
}

type AuthResponse struct {
	AccessToken string      `json:"access_token,omitempty"`
	User        *users.User `json:"user,omitempty"`
}

type Config struct {
	Logger       *zap.Logger
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
		http:          httpClient,
		github:        github.NewClient(httpClient),
		users:         cfg.UserService,
		tokenizer:     cfg.Tokenizer,
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
	sub.HandleFunc("/github", s.AuthHandler)
	sub.HandleFunc("/authorize", s.AuthorizeHandler)
	sub.HandleFunc("/health", s.tokenizer.Authenticator(http.HandlerFunc(s.HealthHandler)).ServeHTTP)
	sub.HandleFunc("/redirect", s.CallbackHandler)
	sub.HandleFunc("", s.AuthenticateHandler)

	sub.ServeHTTP(w, r)
}

func (s *Service) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "auth/github.AuthenticateHandler")
	defer span.End()

	var req auth.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := s.Authenticate(ctx, &req)
	switch {
	case err == nil:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.log.Error("failed to encode auth response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case err == users.ErrNotAuthorized:
		w.WriteHeader(http.StatusUnauthorized)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) Authenticate(ctx context.Context, r *auth.AuthRequest) (*auth.AuthResponse, error) {
	u, err := s.users.Get(ctx, &users.GetParams{
		Email: r.GetEmail(),
	})
	if err != nil {
		return nil, err
	}

	if !users.ComparePasswords(
		[]byte(u.Password),
		[]byte(r.GetPassword()),
	) {
		return nil, users.ErrNotAuthorized
	}
	u.Password = ""

	token, err := s.tokenizer.GenerateToken(ctx, u)
	if err != nil {
		return nil, err
	}

	return &auth.AuthResponse{
		User:  u,
		Token: token,
	}, nil
}

func (s *Service) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		s.log.Error("failed to retrieve auth token", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = cookie.Value
}

func (s *Service) AuthHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("redirecting to github oauth endpoint")

	http.Redirect(
		w,
		r,
		fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", s.clientID, s.redirectURL),
		http.StatusTemporaryRedirect,
	)
}

func (s *Service) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("handling authentication callback")

	// First, we need to get the value of the `code` query param
	err := r.ParseForm()
	if err != nil {
		s.log.Error("failed to parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	code := r.FormValue("code")

	// Next, lets for the HTTP request to call the github oauth enpoint
	// to get our access token
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
			s.clientID,
			s.clientSecret,
			code,
		),
		nil,
	)
	if err != nil {
		s.log.Error("failed to create HTTP request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// We set this header since we want the response
	// as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := s.http.Do(req)
	if err != nil {
		s.log.Error("failed to execute HTTP request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var t AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		s.log.Error("failed to encode body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	client := github.NewClient(
		oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: t.AccessToken,
			},
		)),
	)

	user, _, err := client.Users.Get(r.Context(), "")
	if err != nil {
		s.log.Error("failed to fetch user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	authUser, err := s.users.Get(r.Context(), &users.GetParams{
		AuthenticatedAccounts: []*users.AuthenticatedAccount{
			{
				Id:   strconv.Itoa(int(user.GetID())),
				Type: "github",
			},
		},
	})
	switch {
	case err == nil:
	case err == users.ErrNotFound:
		now := &timestamp.Timestamp{
			Seconds: time.Now().UTC().Unix(),
		}
		authUser, err = s.users.Save(r.Context(), &users.User{
			Id:      uuid.New().String(),
			Name:    user.GetName(),
			Email:   user.GetEmail(),
			Created: now,
			Updated: now,
			AuthenticatedAccounts: []*users.AuthenticatedAccount{
				{
					Id:   strconv.Itoa(int(user.GetID())),
					Type: "github",
				},
			},
		})
		if err != nil {
			s.log.Error("failed to save user record", zap.Error(err))
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	default:
		s.log.Error("failed to fetch user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	token, err := s.tokenizer.GenerateToken(r.Context(), authUser)
	if err != nil {
		s.log.Error("failed to generate auth token", zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Finally, send a response to redirect the user to the "welcome" page
	// with the access token
	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieName,
		Value:    token,
		Domain:   r.URL.Host,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Secure:   true,
		HttpOnly: true,
	})

	s.log.Debug("writing token to response body")

	w.WriteHeader(http.StatusFound)
}

func (s *Service) HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.Health(r.Context(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) Health(_ context.Context, _ *auth.HealthRequest) (*auth.HealthResponse, error) {
	return &auth.HealthResponse{
		Uptime:  time.Since(s.initTimestamp).String(),
		Machine: s.hostname,
	}, nil
}
