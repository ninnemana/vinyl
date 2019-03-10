package rest

import (
	"context"
	lg "log"
	"mime"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ninnemana/vinyl/pkg/log"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func Start(ctx context.Context, serviceAddr, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := log.Init(); err != nil {
		lg.Fatalf("failed to create logger: %v", err)
	}

	gwmux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mime.AddExtensionType(".svg", "image/svg+xml")
	mux.Handle("/openapi/", http.StripPrefix("/openapi", http.FileServer(http.Dir("openapi"))))

	err := vinyl.RegisterVinylHandlerFromEndpoint(
		ctx,
		gwmux,
		serviceAddr,
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		return errors.Wrap(err, "failed to register vinyl handler")
	}

	srv := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	lg.Printf("starting HTTP server on '%s'\n", httpPort)
	return errors.Wrap(srv.ListenAndServe(), "fell out of serving HTTP traffic")
}
