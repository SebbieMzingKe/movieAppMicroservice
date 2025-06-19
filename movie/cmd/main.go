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
	grpchandler "movieapp.com/movie/internal/handler/grpc"
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

	svc := movie.New(ratingGateway, metadataGateway)
	h := grpchandler.New(svc) 
	lis, err := net.Listen("tcp", "localhost:8083")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	srv.Serve(lis)
	// http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}

// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"net"
// 	"time"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection" // This is correctly imported and used
// 	"movieapp.com/gen"
// 	"movieapp.com/movie/internal/controller/movie"
// 	metadatagateway "movieapp.com/movie/internal/gateway/metadata/http"
// 	grpchandler "movieapp.com/movie/internal/handler/grpc"
// 	ratinggateway "movieapp.com/movie/internal/gateway/rating/http"
// 	"movieapp.com/pkg/discovery"
// 	"movieapp.com/pkg/discovery/consul"
// )

// const serviceName = "movie"

// func main() {
// 	var port int

// 	flag.IntVar(&port, "port", 8083, "API handler port")
// 	flag.Parse()

// 	log.Printf("Starting the movie service on port %d", port)

// 	registry, err := consul.NewRegistry("localhost:8500")
// 	if err != nil {
// 		log.Fatalf("Failed to create registry: %v", err) // Use log.Fatalf for critical startup errors
// 	}

// 	ctx := context.Background()

// 	instanceID := discovery.GenerateInstanceID(serviceName)

// 	// Register the service with Consul
// 	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
// 		log.Fatalf("Failed to register service: %v", err)
// 	}

// 	// Ensure deregistration on exit, regardless of how the program exits
// 	// Only one defer is needed here
// 	defer func() {
// 		log.Printf("Deregistering service %s/%s", serviceName, instanceID)
// 		if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
// 			log.Printf("Failed to deregister service: %v", err)
// 		}
// 	}()

// 	// Goroutine to report healthy state to Consul
// 	go func() {
// 		for {
// 			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
// 				log.Println("Failed to report healthy state: " + err.Error())
// 			}
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()

// 	metadataGateway := metadatagateway.New(registry)
// 	ratingGateway := ratinggateway.New(registry)
// 	svc := movie.New(ratingGateway, metadataGateway)
// 	h := grpchandler.New(svc)

// 	// Listen for TCP connections
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) // Use ":%d" for listening on all interfaces
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}

// 	// Create a new gRPC server
// 	srv := grpc.NewServer()

// 	// Register the reflection service (this is correct!)
// 	reflection.Register(srv)

// 	// Register your MovieService implementation
// 	gen.RegisterMovieServiceServer(srv, h)

// 	log.Printf("gRPC server serving on %s", lis.Addr())
// 	// Start serving gRPC requests (this call is blocking)
// 	if err := srv.Serve(lis); err != nil {
// 		// If srv.Serve returns an error, it means the server stopped unexpectedly
// 		log.Fatalf("Failed to serve: %v", err)
// 	}
// }