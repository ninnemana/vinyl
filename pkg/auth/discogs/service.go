package discogs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	httpserver "github.com/ninnemana/vinyl/pkg/http"

	"go.opentelemetry.io/otel"

	"github.com/gorilla/mux"
	"github.com/ninnemana/tracelog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	Entity = "discogs"
)

var (
	tracer = otel.Tracer("pkg/auth/discogs")
)

type Service struct {
	log *tracelog.TraceLogger
}

func New(ctx context.Context, l *tracelog.TraceLogger) (*Service, error) {
	ctx, span := tracer.Start(ctx, "discogs/service.New")
	defer span.End()

	svc := &Service{
		l,
	}

	if err := httpserver.RegisterHandler(svc); err != nil {
		return nil, fmt.Errorf("failed to register discogs handler with HTTP server: %w", err)
	}

	return svc, nil
}

func (s *Service) Middleware() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}

func (s *Service) Register(rpc *grpc.Server) error {
	s.log.Debug(
		"register RPC service",
		zap.Any("info", rpc.GetServiceInfo()),
	)

	RegisterDiscogsServer(rpc, s)

	return nil
}

func (s *Service) Route() string {
	return "/" + Entity
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "discogs/service.ServeHTTP")
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
	sub.Methods(http.MethodGet).Path("").HandlerFunc(s.RequestTokenHandler)
	sub.Methods(http.MethodPost).Path("").HandlerFunc(s.AccessTokenHandler)

	sub.ServeHTTP(w, r)
}

func (s *Service) RequestTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "discogs/service.RequestTokenHandler")
	defer span.End()

	result, err := s.RequestToken(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "discogs/service.AccessTokenHandler")
	defer span.End()

	var params AccessTokenParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	result, err := s.AccessToken(ctx, &params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) RequestToken(ctx context.Context, params *RequestTokenParams) (*RequestTokenResult, error) {
	ctx, span := tracer.Start(ctx, "discogs/service.RequestToken")
	defer span.End()

	req, err := http.NewRequest(http.MethodGet, "https://api.discogs.com/oauth/request_token", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization",
		fmt.Sprintf("OAuth %s",
			strings.Join([]string{
				fmt.Sprintf("oauth_nonce=%d", time.Now().UTC().UnixNano()),
				fmt.Sprintf("oauth_signature_method=%s", "PLAINTEXT"),
				fmt.Sprintf("oauth_timestamp=%d", time.Now().UTC().Unix()),
				fmt.Sprintf("oauth_callback=%s", "/discogs/callback"),
				fmt.Sprintf("oauth_signature=%s&", "CBwRadHPnhlMzLmigCDbeMPSxEsKibfe"),
				fmt.Sprintf("oauth_consumer_key=%s", "bfiOLMABZAdENUAdaLPx"),
			}, ","),
		),
	)
	req.Header.Set("User-Agent", "go-client-vinyl")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to retrieve a request token was not valid: %s", resp.Status)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("the returned token could not be decoded: %w", err)
	}

	return &RequestTokenResult{
		Token: result.Get("oauth_token"),
	}, nil
}

func (s *Service) AccessToken(ctx context.Context, params *AccessTokenParams) (*AccessTokenResult, error) {
	ctx, span := tracer.Start(ctx, "discogs/service.AccessToken")
	defer span.End()

	req, err := http.NewRequest(http.MethodGet, "https://api.discogs.com/oauth/access_token", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization",
		fmt.Sprintf("OAuth %s",
			strings.Join([]string{
				fmt.Sprintf("oauth_token=%s", params.GetRequestToken()),
				fmt.Sprintf("oauth_nonce=%d", time.Now().UTC().UnixNano()),
				fmt.Sprintf("oauth_signature_method=%s", "PLAINTEXT"),
				fmt.Sprintf("oauth_timestamp=%d", time.Now().UTC().Unix()),
				fmt.Sprintf("oauth_callback=%s", "/discogs/callback"),
				fmt.Sprintf("oauth_signature=%s&", "CBwRadHPnhlMzLmigCDbeMPSxEsKibfe"),
				fmt.Sprintf("oauth_consumer_key=%s", "bfiOLMABZAdENUAdaLPx"),
			}, ","),
		),
	)
	req.Header.Set("User-Agent", "go-client-vinyl")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to retrieve a request token was not valid: %s", resp.Status)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("the returned token could not be decoded: %w", err)
	}

	return &AccessTokenResult{
		Token:  result.Get("oauth_token"),
		Secret: result.Get("oauth_token_secret"),
	}, nil
}

func (s *Service) mustEmbedUnimplementedDiscogsServer() {
	//TODO implement me
	panic("implement me")
}
