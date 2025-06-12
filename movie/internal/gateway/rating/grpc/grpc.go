package grpc

import (
	"context"

	"movieapp.com/gen"
	"movieapp.com/internal/grpcutil"
	"movieapp.com/pkg/discovery"
	model "movieapp.com/rating/pkg"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordId, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)

	if err != nil {
		return 0, err
	}

	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)})

	if err != nil {
		return 0, err
	}

	return resp.RatingValue, nil
}