package router

import (
	"context"
	"fmt"
	"os"

	httpserver "github.com/ninnemana/vinyl/pkg/http"
	"github.com/ninnemana/vinyl/pkg/vinyl/firestore"
	"go.uber.org/zap"
)

func Initialize(log *zap.Logger) error {
	svc, err := firestore.New(
		context.Background(),
		log,
		os.Getenv("GCP_PROJECT_ID"),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create firestore service: %w", err)
	}

	if err := httpserver.RegisterHandler(svc); err != nil {
		return err
	}

	return nil
}
