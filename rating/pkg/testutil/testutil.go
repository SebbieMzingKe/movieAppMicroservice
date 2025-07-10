package testutil

import (
	"movieapp.com/gen"
	"movieapp.com/rating/internal/handler/grpc"
	"movieapp.com/rating/internal/controller/rating"
	"movieapp.com/rating/internal/repository/memory"
)

func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpc.New(ctrl)
}