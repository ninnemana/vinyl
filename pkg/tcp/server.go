package tcp

import (
	"context"

	"github.com/ninnemana/drudge"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"github.com/ninnemana/vinyl/pkg/vinyl/firestore"
	"github.com/uber/jaeger-client-go/config"
)

const (
	tcpAddr = ":8080"
	rpcAddr = ":8081"
)

var (
	options = drudge.Options{
		BasePath: "/",
		Addr:     tcpAddr,
		RPC: drudge.Endpoint{
			Network: "tcp",
			Addr:    rpcAddr,
		},
		SwaggerDir:    "openapi",
		Handlers:      []drudge.Handler{vinyl.RegisterVinylHandler},
		OnRegister:    firestore.Register,
		TraceExporter: drudge.Jaeger,
	}
)

func Serve() error {
	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}

	cfg.ServiceName = "vinyltap"
	options.TraceConfig = cfg

	return drudge.Run(context.Background(), options)
}
