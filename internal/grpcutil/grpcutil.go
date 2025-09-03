package grpcutil

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"movieapp.com/pkg/discovery"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// ServiceConnection trys to selects a random service instance and returns a grpc connection to it
func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(
		addrs[rand.Intn(len(addrs))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
}

