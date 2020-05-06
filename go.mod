module github.com/ninnemana/vinyl

go 1.13

require (
	cloud.google.com/go v0.52.0 // indirect
	cloud.google.com/go/firestore v1.1.1
	contrib.go.opencensus.io/exporter/jaeger v0.2.0 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.7.4
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.11.3
	github.com/irlndts/go-discogs v0.2.5
	github.com/ninnemana/drudge v0.1.2-0.20200328191329-a1c1b087f750
	github.com/opentracing/opentracing-go v1.1.1-0.20190913142402-a7454ce5950e
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.4.0 // indirect
	github.com/soheilhy/cmux v0.1.4
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	go.opencensus.io v0.22.2
	go.opentelemetry.io/otel v0.3.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.3.0
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0
	golang.org/x/exp v0.0.0-20200119233911-0405dc783f0a // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	golang.org/x/tools v0.0.0-20200130002326-2f3ba24bd6e7 // indirect
	google.golang.org/api v0.20.0
	google.golang.org/genproto v0.0.0-20200128133413-58ce757ed39b
	google.golang.org/grpc v1.27.1
	k8s.io/apimachinery v0.18.2
)

replace github.com/go-logfmt/logfmt => github.com/go-logfmt/logfmt v0.3.0
