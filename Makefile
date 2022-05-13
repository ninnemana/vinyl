PKG="${GOPATH}/src/github.com/ninnemana/vinyl"
SHELL := /bin/bash -o pipefail

# proto is a target that uses prototool.
# By depending on $(PROTOTOOL), prototool will be installed on the Makefile's path.
# Since the path above has the temporary GOBIN at the front, this will use the
# locally installed prototool.
.PHONY: generate

godeps:
	@go get github.com/grpc-ecosystem/grpc-gateway@v1.14.5
	@go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.5
	@go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.14.5
	@go get google.golang.org/protobuf@v1.28.0
	@go get github.com/gogo/protobuf/gogoproto@v1.3.1
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0


protoc:
	@./scripts/protoc.sh

generate: godeps
	@go generate ./...

build: generate
	go build -v ./cmd/server
	go build -v ./cmd/client

test: generate
	go test -cover ./...

gen_cert:
	rm -rf certs
	mkdir -p certs
	openssl genrsa -out ./certs/server.key 2048
	openssl req -new -x509 -key ./certs/server.key -out ./certs/server.pem -days 3650

run: generate
	@go run ./cmd/server

build-ui:
	@cd ui; npm run-script build; cd ../

gen_docs: generate
	@npm install -g redoc-cli
	@redoc-cli bundle \
		${PKG}/openapi/vinyl.swagger.json \
		-o="${PKG}/openapi/index.html" \
		--title "Vinyl Registry API"

.PHONY: generate
