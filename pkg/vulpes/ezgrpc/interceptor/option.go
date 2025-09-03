// Package interceptor provides gRPC unary server interceptors for common concerns
// such as logging, metrics, rate limiting, and panic recovery.
package interceptor

import (
	"google.golang.org/grpc"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

// interceptors is a slice of gRPC unary server interceptors that are chained together.
// The order of interceptors is important as they are executed sequentially.
var interceptors []grpc.UnaryServerInterceptor = []grpc.UnaryServerInterceptor{
	// 1. Recovery: The outermost interceptor to catch any panics from downstream handlers.
	recoveryInterceptor,

	// 2. Prometheus: Provides monitoring metrics for gRPC requests.
	grpc_prometheus.UnaryServerInterceptor,

	// 3. RequestID: Ensures each request has a unique identifier.
	requestIDInterceptor,

	// 4. Logger: Logs detailed information about each request, depends on RequestID.
	loggerInterceptor,

	// 5. RateLimit: Rejects requests early to save resources.
	rateLimitInterceptor,

	// 6. Validation: The last interceptor to run, ensuring that only valid requests are processed.
	validateUnaryInterceptor,
}

// NewGrpcServerWithInterceptors creates a new gRPC server with the predefined chain of unary interceptors.
// This simplifies server setup by providing a standard set of middleware.
func NewGrpcServerWithInterceptors() *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors...,
		),
	)
}
