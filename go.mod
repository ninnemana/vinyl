module github.com/ninnemana/vinyl

go 1.13

require (
	cloud.google.com/go v0.69.1 // indirect
	cloud.google.com/go/firestore v1.3.0
	cloud.google.com/go/logging v1.1.0
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/stackdriver v0.13.4
	github.com/caarlos0/env/v6 v6.3.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.4.2
	github.com/google/go-github/v31 v31.0.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/grpc-gateway v1.15.2
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/ninnemana/go-discogs v0.2.5-0.20200606155508-e2bcc3acd5a9
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/soheilhy/cmux v0.1.4
	go.opencensus.io v0.22.5
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/sys v0.0.0-20201015000850-e3ed0017c211 // indirect
	google.golang.org/api v0.33.0
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201014134559-03b6142f0dc9
	google.golang.org/grpc v1.33.0
	google.golang.org/grpc/examples v0.0.0-20201014215113-7b167fd6eca1 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
)

replace github.com/go-logfmt/logfmt => github.com/go-logfmt/logfmt v0.3.0
