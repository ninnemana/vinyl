package auth

//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative ./auth.proto --proto_path=$GOPATH/src/github.com/ninnemana/vinyl --proto_path=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/ --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/ --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel"

	"github.com/dgrijalva/jwt-go"
	"github.com/ninnemana/vinyl/pkg/users"
)

const (
	CookieName = "vinyl-auth"
)

var (
	tracer = otel.Tracer("pkg/auth")
)

type User struct {
	User  *users.User `json:"user"`
	Token string      `json:"token"`
}

type UserClaims struct {
	Authorized bool   `json:"authorized,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	Expires    int64  `json:"expires,omitempty"`
	Token      string `json:"token,omitempty"`
}

func (u UserClaims) Valid() error {
	return nil
}

type Tokenizer interface {
	GenerateToken(context.Context, *users.User) (string, error)
	Authenticator(http.Handler) http.Handler
}

type JWT struct {
	accessSecret string
	replacer     *strings.Replacer
	validator    func(*jwt.Token) (interface{}, error)
}

func NewTokenizer(accessSecret string) (Tokenizer, error) {
	if accessSecret == "" {
		return nil, errors.New("no JWT access token was provided")
	}

	return &JWT{
		accessSecret: accessSecret,
		replacer:     strings.NewReplacer("Bearer ", "", "bearer", ""),
		validator: func(token *jwt.Token) (interface{}, error) {
			return []byte(accessSecret), nil
		},
	}, nil
}

func (t *JWT) GenerateToken(_ context.Context, u *users.User) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Authorized: true,
		UserID:     u.GetId(),
		Expires:    time.Now().Add(time.Minute * 15).Unix(),
	})

	token, err := at.SignedString([]byte(t.accessSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (t *JWT) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "jwt.Authenticator")
		defer span.End()

		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Not Authorized"))

			span.SetStatus(codes.Ok, "Not Authorized")
			span.AddEvent("unauthorized")

			return
		}

		token, err := jwt.ParseWithClaims(
			t.replacer.Replace(auth),
			&UserClaims{},
			t.validator,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		c, ok := token.Claims.(*UserClaims)
		if !ok {
			http.Error(w, "claim was invalid", http.StatusUnauthorized)
			return
		}
		_ = c

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
