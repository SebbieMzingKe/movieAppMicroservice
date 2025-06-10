package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"movieapp.com/gen"
	"movieapp.com/metadata/internal/controller/metadata"
	grpchandler "movieapp.com/metadata/internal/handler/grpc"
	"movieapp.com/metadata/internal/repository/memory"
	"movieapp.com/pkg/discovery"
	"movieapp.com/pkg/discovery/consul"
)
const serviceName = "metadata"

func main() {

	var port int

	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()

	log.Printf("Starting the metadata service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")

	if err != nil {
		panic(err)
	}
	
	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {

				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	log.Println("Starting the movie metadata service")
	repo := memory.New()
	svc := metadata.New(repo)
	h := grpchandler.New(svc)

	lis, err := net.Listen("tcp", "localhost:8-81")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	srv.Serve(lis)
	
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}