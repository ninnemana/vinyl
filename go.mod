module github.com/ninnemana/vinyl

go 1.13

require (
	cloud.google.com/go/firestore v1.1.1
	cloud.google.com/go/logging v1.1.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.12.0
	github.com/blendle/zapdriver v1.3.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/google/go-github/v31 v31.0.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.11.3
	github.com/jonstaryuk/gcloudzap v0.1.1
	github.com/ninnemana/go-discogs v0.2.5-0.20200606155508-e2bcc3acd5a9
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/soheilhy/cmux v0.1.4
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200827165113-ac2560b5e952
	google.golang.org/grpc v1.32.0
	google.golang.org/grpc/examples v0.0.0-20201010204749-3c400e7fcc87 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/go-logfmt/logfmt => github.com/go-logfmt/logfmt v0.3.0
