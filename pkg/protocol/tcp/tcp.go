package tcp

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ninnemana/drudge"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"

	"github.com/ninnemana/vinyl/pkg/vinyl"
	vinylService "github.com/ninnemana/vinyl/pkg/vinyl/datastore"
)

const (
	tcpAddr = "vinyltap:8000"
	rpcAddr = "vinyltap:8080"
)

func Serve() error {

	return errors.Wrap(drudge.Run(context.Background(), drudge.Options{
		Metrics: &drudge.Metrics{
			Prefix:      "vinyl",
			PullAddress: ":9090",
		},
		BasePath: "/",
		Addr:     tcpAddr,
		RPC: drudge.Endpoint{
			Network: "tcp",
			Addr:    rpcAddr,
		},
		Handlers:   []drudge.Handler{vinyl.RegisterVinylHandler},
		SwaggerDir: "/openapi",
		Mux: []runtime.ServeMuxOption{
			// we take the JWT token off the HTTP header and place it in
			// the metadata within the context.
			runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
				in, _ := metadata.FromIncomingContext(ctx)
				out, _ := metadata.FromOutgoingContext(ctx)

				return metadata.Join(in, out)
			}),
		},
		OnRegister: vinylService.Register,
	}), "failed to start application server")
}
