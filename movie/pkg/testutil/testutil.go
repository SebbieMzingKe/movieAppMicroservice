package testutil

import (
	"movieapp.com/gen"
	"movieapp.com/movie/internal/controller/movie"
	metadatagateway "movieapp.com/movie/internal/gateway/metadata/grpc"
	ratinggateway "movieapp.com/movie/internal/gateway/rating/grpc"
	grpchandler "movieapp.com/movie/internal/handler/grpc"
	"movieapp.com/pkg/discovery"
)

func NewTestMovieGRPCServer(registry discovery.Registry) gen.MovieServiceServer {
	metadataGateway := metadatagateway.New(registry)
	ratinggateway := ratinggateway.New(registry)
	ctrl := movie.New(ratinggateway, metadataGateway)

	return grpchandler.New(ctrl)
}
