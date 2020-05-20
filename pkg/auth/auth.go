package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ninnemana/vinyl/pkg/users"
)

const (
	CookieName = "vinyl-auth"
)

var (
	accessSecret = os.Getenv("JWT_ACCESS_SECRET")
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

func GenerateToken(u *users.User) (string, error) {
	if accessSecret == "" {
		return "", errors.New("no JWT access token was provided")
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Authorized: true,
		UserID:     u.GetId(),
		Expires:    time.Now().Add(time.Minute * 15).Unix(),
		// Token:      u.Token,
	})

	token, err := at.SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func Authenticator(next http.Handler) http.Handler {
	replacer := strings.NewReplacer("Bearer ", "", "bearer", "")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(replacer.Replace(auth), &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(accessSecret), nil
		})
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

		next.ServeHTTP(w, r)
	})
}
