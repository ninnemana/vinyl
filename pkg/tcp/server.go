package tcp

import (
	"context"

	"github.com/ninnemana/drudge"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	// "github.com/ninnemana/vinyl/pkg/vinyl/datastore"
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
	// options := []option.ClientOption{}

	// if os.Getenv("GCP_SVC_ACCOUNT") != "" {
	// 	js, err := base64.StdEncoding.DecodeString(os.Getenv("GCP_SVC_ACCOUNT"))
	// 	if err != nil {
	// 		return errors.WithMessage(err, "failed to decode service account")
	// 	}

	// 	options = append(options, option.WithCredentialsJSON(js))
	// }

	// client, err := firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT_ID"), options...)
	// if err != nil {
	// 	return errors.WithMessage(err, "failed to create Firestore client")
	// }

	// service := &Service{
	// 	albums: map[int32]*vinyltap.Album{},
	// 	client: client,
	// }

	// vinyltap.RegisterTapServer(server, service)

	return nil
}
