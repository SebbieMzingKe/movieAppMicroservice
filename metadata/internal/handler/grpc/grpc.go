package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieapp.com/gen"
	"movieapp.com/metadata/internal/controller/metadata"
	"movieapp.com/metadata/pkg/model"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	svc *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{svc: ctrl}
}

func (h *Handler) GetMetadataByID(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req orempty id")
	}

	m, err := h.svc.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, metadata.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil	
}


