package tcp

import (
	"context"
	"os"

	"github.com/ninnemana/drudge"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"github.com/ninnemana/vinyl/pkg/vinyl/datastore"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
		OnRegister:    Register,
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

func Register(server *grpc.Server) error {

	log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	service, err := datastore.New(context.Background(), log, os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		return err
	}

	vinyl.RegisterVinylServer(server, service)

	return nil
}
