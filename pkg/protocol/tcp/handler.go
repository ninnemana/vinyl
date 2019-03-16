package tcp

import (
	"net/http"
	"strings"

	"go.opencensus.io/trace"

	"go.uber.org/zap"
)

type Handler struct {
	server *Server
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "serve")
	switch {
	case r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc"):
		span.AddAttributes(trace.StringAttribute("protocol", "grpc"))
		h.rpc(w, r.WithContext(ctx))
	case strings.HasPrefix(r.URL.Path, "/openapi/"):
		span.AddAttributes(trace.StringAttribute("protocol", "openapi"))
		h.openapi(w, r.WithContext(ctx))
	default:
		span.AddAttributes(trace.StringAttribute("protocol", "http"))
		h.http(w, r.WithContext(ctx))
	}
	span.End()
}

func (h *Handler) rpc(w http.ResponseWriter, r *http.Request) {
	h.server.log.Info("Handling RPC Call", zap.String("path", r.URL.Path))
	h.server.grpc.ServeHTTP(w, r)
}

func (h *Handler) http(w http.ResponseWriter, r *http.Request) {
	h.server.log.Info("Handling HTTP Call", zap.String("path", r.URL.Path))
	h.server.gateway.ServeHTTP(w, r)
}

func (h *Handler) openapi(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix(
		"/openapi/",
		http.FileServer(http.Dir("openapi")),
	).ServeHTTP(w, r)
}
