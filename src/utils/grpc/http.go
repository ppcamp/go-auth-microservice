package grpc

import (
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func httpAndGrpcMux(httpHandler http.Handler, grpcHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("content-type"), "application/grpc") {
			grpcHandler.ServeHTTP(w, r)
			return
		}
		httpHandler.ServeHTTP(w, r)
	})
}

// NewMuxServer implements a gRPC http2 connection server using the httpAndGrpcMux to check if
func NewMuxServer(httpHandler http.Handler, grpcHandler http.Handler) *http.Server {
	muxHandlers := httpAndGrpcMux(httpHandler, grpcHandler)
	http2Server := new(http2.Server)

	return &http.Server{Handler: h2c.NewHandler(muxHandlers, http2Server)}
}
