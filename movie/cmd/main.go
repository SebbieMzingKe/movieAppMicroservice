package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"movieapp.com/gen"
	"movieapp.com/movie/internal/controller/movie"
	metadatagateway "movieapp.com/movie/internal/gateway/metadata/http"
	grpchandler "movieapp.com/movie/internal/gateway/rating/grpc"
	ratinggateway "movieapp.com/movie/internal/gateway/rating/http"
	"movieapp.com/pkg/discovery"
	"movieapp.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {

	var port int

	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()

	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")

	if err != nil {
		panic(err)
	}

	// registry := consul.NewRegistry(map[string][]string{
	// 	"metadata": {"localhost:8081"},
	// 	"rating": {"localhost:8082"},
	// 	"movie": {"localhost:8083"},
	// })
	
	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	defer registry.Deregister(ctx, instanceID, serviceName)

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {

				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)

	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	// http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}