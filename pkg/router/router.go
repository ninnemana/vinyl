package router

import (
	"context"
	"fmt"
	"os"

	"github.com/ninnemana/vinyl/pkg/auth/github"
	httpserver "github.com/ninnemana/vinyl/pkg/http"
	userStore "github.com/ninnemana/vinyl/pkg/users/firestore"
	vinylStore "github.com/ninnemana/vinyl/pkg/vinyl/firestore"
	"go.uber.org/zap"
)

func Initialize(log *zap.Logger) error {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")

	svc, err := vinylStore.New(
		context.Background(),
		log,
		projectID,
		nil,
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
		log,
		userSvc,
	)
	if err != nil {
		return fmt.Errorf("failed to create github service: %w", err)
	}

	if err := httpserver.RegisterHandler(githubSvc); err != nil {
		return err
	}

	return nil
}
