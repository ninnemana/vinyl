package vinyl

//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./vinyl.proto --proto_path=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/ --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/ --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound         = errors.Errorf("no vinyl was found")
	ErrInvalidGetParams = errors.New("invalid parameters supplied when attempting to retrieve a vinyl")
)
