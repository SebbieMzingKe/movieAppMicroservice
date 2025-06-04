package grpc

import (
	"movieapp.com/gen"
	"movieapp.com/metadata/internal/controller/metadata"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}