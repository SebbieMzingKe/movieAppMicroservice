package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieapp.com/gen"
	"movieapp.com/rating/internal/controller/rating"
	model "movieapp.com/rating/pkg"
)

type Handler struct {
	gen.UnimplementedRatingServiceServer
	svc *rating.Controller
}

func New(svc *rating.Controller) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil request or empty id")
	} 

	v, err := h.svc.GetAggregatedRating(ctx, model.RecordId(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}


func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil request or empty user id or record id")
	}

	if err := h.svc.PutRating(ctx, model.RecordId(req.RecordId), model.RecordType(req.RecordType), &model.Rating{UserID: model.UserID(req.UserId), Value: model.RatingValue(req.RecordValue)}); err != nil {
		return nil, err
	}
	return &gen.PutRatingResponse{}, nil
}