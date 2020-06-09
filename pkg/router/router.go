package router

import (
	"context"
	"fmt"
	"os"

	"github.com/ninnemana/vinyl/pkg/auth"
	"github.com/ninnemana/vinyl/pkg/auth/github"
	httpserver "github.com/ninnemana/vinyl/pkg/http"
	userStore "github.com/ninnemana/vinyl/pkg/users/firestore"
	vinylStore "github.com/ninnemana/vinyl/pkg/vinyl/firestore"
	"go.uber.org/zap"
)

var (
	discogsKey   = os.Getenv("DISCOGS_API_KEY")
	projectID    = os.Getenv("GCP_PROJECT_ID")
	jwtSecret    = os.Getenv("JWT_ACCESS_SECRET")
	redirectURL  = os.Getenv("BASE_URL") + "/auth/redirect"
	clientID     = os.Getenv("GITHUB_CLIENT_ID")
	clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
)

func Initialize(log *zap.Logger) error {
	ctx := context.Background()

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to fetch hostname: %w", err)
	}

	tokenizer, err := auth.NewTokenizer(jwtSecret)
	if err != nil {
		return err
	}

	svc, err := vinylStore.New(
		context.Background(),
		vinylStore.Config{
			Logger:          log,
			GoogleProjectID: projectID,
			DiscogsAPIKey:   discogsKey,
			Hostname:        hostname,
			Tokenizer:       tokenizer,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create firestore service: %w", err)
	}

	if err := httpserver.RegisterHandler(svc); err != nil {
		return err
	}

	userSvc, err := userStore.New(ctx, log, projectID)
	if err != nil {
		return err
	}

	if err := httpserver.RegisterHandler(userSvc); err != nil {
		return err
	}

	githubSvc, err := github.New(
		context.Background(),
		github.Config{
			Logger:       log,
			UserService:  userSvc,
			Tokenizer:    tokenizer,
			Hostname:     hostname,
			RedirectURL:  redirectURL,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create github service: %w", err)
	}

	if err := httpserver.RegisterHandler(githubSvc); err != nil {
		return err
	}

	return nil
}
