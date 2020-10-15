package vinyl

//go:generate protoc --include_imports --include_source_info --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.5/ --proto_path=$GOPATH/pkg/mod/ --proto_path=$GOPATH/src/github.com/ninnemana/vinyl/ --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.5/third_party/googleapis/ --proto_path=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.1/ --descriptor_set_out=api_descriptor.pb --go_out=plugins=grpc:$GOPATH/src $GOPATH/src/github.com/ninnemana/vinyl/pkg/vinyl/vinyl.proto

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound         = errors.Errorf("no vinyl was found")
	ErrInvalidGetParams = errors.New("invalid parameters supplied when attempting to retrieve a vinyl")
)
