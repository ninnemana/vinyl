package discogs

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/gomodule/oauth1/oauth"
	"github.com/gorilla/mux"
	"github.com/ninnemana/go-discogs"
	"github.com/ninnemana/tracelog"
	"github.com/ninnemana/vinyl/pkg/auth"
	"github.com/ninnemana/vinyl/pkg/users"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	Entity = "auth"
)

var (
	ErrInvalidLogger = errors.New("the provided logger was not valid")

	oauthClient = oauth.Client{
		TemporaryCredentialRequestURI: "https://api.discogs.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://www.discogs.com/oauth/authorize",
		TokenRequestURI:               "https://api.discogs.com/oauth/access_token",
		Header:                        http.Header{"User-Agent": {"ExampleDiscogsClient/1.0"}},
	}

	tmpCreds *oauth.Credentials
	tracer   = otel.Tracer("pkg/auth/discogs")
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
	Logger        *tracelog.TraceLogger
	UserService   users.UsersServer
	Tokenizer     auth.Tokenizer
	DiscogsAPIKey string
	Hostname      string
	RedirectURL   string
	ClientID      string
	ClientSecret  string
}

func New(ctx context.Context, cfg Config) (*Service, error) {
	if cfg.Logger == nil {
		return nil, ErrInvalidLogger
	}

	disc, err := discogs.New(&discogs.Options{
		UserAgent: "Some Agent",
		Token:     cfg.DiscogsAPIKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create discogs client")
	}

	return &Service{
		log:           cfg.Logger,
		hostname:      cfg.Hostname,
		initTimestamp: time.Now().UTC(),
		users:         cfg.UserService,
		redirectURL:   cfg.RedirectURL,
		clientID:      cfg.ClientID,
		clientSecret:  cfg.ClientSecret,
		discogs:       disc,
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
	// sub.HandleFunc("/github", s.AuthHandler)
	// sub.HandleFunc("/authorize", s.AuthorizeHandler)
	// sub.HandleFunc("/health", s.tokenizer.Authenticator(http.HandlerFunc(s.HealthHandler)).ServeHTTP)
	sub.HandleFunc("/token", s.TokenHandler)
	sub.HandleFunc("/callback", s.CallbackHandler)
	// sub.HandleFunc("", s.AuthenticateHandler)

	sub.ServeHTTP(w, r)
}

func (s *Service) Authenticate(ctx context.Context, r *auth.AuthRequest) (*auth.AuthResponse, error) {
	return nil, nil
}

func (s *Service) Callback(ctx context.Context, r *auth.CallbackRequest) (*auth.CallbackResponse, error) {
	return nil, nil
}

func (s *Service) Health(ctx context.Context, r *auth.HealthRequest) (*auth.HealthResponse, error) {
	return nil, nil
}

func (s *Service) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "auth/discogs.Callback")
	defer span.End()
	_ = ctx

	token := r.FormValue("oauth_token")
	verifier := r.FormValue("oauth_verifier")

	s.log.Debug(
		"handling Discogs OAuth callback",
		zap.String("token", token),
		zap.String("verifier", verifier),
	)

	tokenCred, _, err := oauthClient.RequestToken(
		nil,
		tmpCreds,
		verifier,
	)
	if err != nil {
		http.Error(w, "Error getting request token, "+err.Error(), 500)
		return
	}

	id, err := s.discogs.OAuthIdentity(
		ctx,
		discogs.WithClient(&oauthClient),
		discogs.WithCredentials(tokenCred),
	)
	if err != nil {
		http.Error(w, "failed to fetch identity: "+err.Error(), 500)
		return
	}

	coll, err := s.discogs.GetFolders(ctx, id.Username, discogs.WithCredentials(tokenCred), discogs.WithClient(&oauthClient))
	if err != nil {
		http.Error(w, "failed to fetch folders: "+err.Error(), 500)
		return
	}

	data, err := json.Marshal(coll)
	if err != nil {
		http.Error(w, "failed to fetch folders: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		http.Error(w, "failed to fetch folders: "+err.Error(), 500)
		return
	}
}

func (s *Service) TokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "auth/github.DiscogsToken")
	defer span.End()

	s.log.Debug("generating Discogs OAuth token")

	oauthClient.Credentials = oauth.Credentials{
		Token:  "bfiOLMABZAdENUAdaLPx",
		Secret: "CBwRadHPnhlMzLmigCDbeMPSxEsKibfe",
	}

	// Next, lets for the HTTP request to call the github oauth enpoint
	// to get our access token
	callback := "http://" + r.Host + "/auth/discogs/callback"
	tempCred, err := oauthClient.RequestTemporaryCredentialsContext(ctx, callback, nil)
	if err != nil {
		http.Error(w, "Error getting temp cred, "+err.Error(), 500)
		return
	}

	tmpCreds = tempCred

	http.Redirect(w, r.WithContext(ctx), oauthClient.AuthorizationURL(tempCred, nil), 302)
}
