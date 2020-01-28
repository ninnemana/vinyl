PKG="${GOPATH}/src/github.com/ninnemana/vinyl"
SHELL := /bin/bash -o pipefail

UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)


# proto is a target that uses prototool.
# By depending on $(PROTOTOOL), prototool will be installed on the Makefile's path.
# Since the path above has the temporary GOBIN at the front, this will use the
# locally installed prototool.
.PHONY: generate

godeps:
	@go get github.com/grpc-ecosystem/grpc-gateway@v1.11.3
	@go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.11.3
	@go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.11.3
	@go get github.com/golang/protobuf/protoc-gen-go@v1.3.2

prototool: godeps
	@curl -sSL \
		https://github.com/uber/prototool/releases/download/v1.9.0/prototool-Darwin-x86_64 \
		-o ./prototool
	@echo "making prototool executable"
	@chmod u+x ./prototool

generate: prototool
	@./prototool generate prototool.yaml --debug
	@npm install -g redoc-cli
	@redoc-cli bundle \
		${PKG}/openapi/vinyl.swagger.json \
		-o="${PKG}/openapi/index.html" \
		--title "Vinyl Registry API"

gen_cert:
	rm -rf certs
	mkdir -p certs
	openssl genrsa -out ./certs/server.key 2048
	openssl req -new -x509 -key ./certs/server.key -out ./certs/server.pem -days 3650

.PHONY: generate