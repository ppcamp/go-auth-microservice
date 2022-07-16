package grpc

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type GrpcHealthService struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (m *GrpcHealthService) Check(_ context.Context, _ *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (m *GrpcHealthService) Watch(_ *grpc_health_v1.HealthCheckRequest, _ grpc_health_v1.Health_WatchServer) error {
	panic("not implemented")
}
