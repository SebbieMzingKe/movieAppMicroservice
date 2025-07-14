package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieapp.com/gen"
	"movieapp.com/internal/grpcutil"
	"movieapp.com/metadata/pkg/model"
	"movieapp.com/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {		
		resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
		if err != nil {
			if shouldRetry(err) {
				continue
			}
			return nil, err
		}
	
		return model.MetadataFromProto(resp.Metadata), nil
	}
	return nil, err

}

func shouldRetry(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == codes.DeadlineExceeded || e.Code() == codes.ResourceExhausted || e.Code() == codes.Unavailable
}