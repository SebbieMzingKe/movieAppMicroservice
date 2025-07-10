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

func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
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

// PutMetadata puts movie metadata to repository.
func (h *Handler) PutMetadata(ctx context.Context, req *gen.PutMetadataRequest) (*gen.PutMetadataResponse, error) {
	if req == nil || req.Metadata == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or metadata")
	}
	metadata := model.MetadataFromProto(req.Metadata)
	if err := h.svc.Put(ctx, metadata); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &gen.PutMetadataResponse{}, nil
}
