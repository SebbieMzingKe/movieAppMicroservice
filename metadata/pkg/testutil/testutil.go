package testutil

import (
	"movieapp.com/gen"
	"movieapp.com/metadata/internal/controller/metadata"
	"movieapp.com/metadata/internal/handler/grpc"
	"movieapp.com/metadata/internal/repository/memory"
)

func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpc.New(ctrl)
}