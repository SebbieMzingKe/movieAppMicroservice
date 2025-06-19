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
	"google.golang.org/grpc/reflection"
	"movieapp.com/gen"
	"movieapp.com/pkg/discovery"
	"movieapp.com/pkg/discovery/consul"
	"movieapp.com/rating/internal/controller/rating"
	grpchandler "movieapp.com/rating/internal/handler/grpc"

	// "movieapp.com/rating/internal/ingester/kafka"
	"movieapp.com/rating/internal/repository/mysql"
)

const serviceName = "rating"

func main() {

	var port int

	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()

	log.Printf("Starting the rating service on port %d", port)
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

	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}
	// ingester, err := kafka.NewIngester("localhost", "rating", "ratings")
	// if err != nil {
	// 	log.Fatalf("failed to initialize ingester: %v", err)
	// }

	svc := rating.New(repo, nil)
	h := grpchandler.New(svc)

	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	reflection.Register(srv)
	srv.Serve(lis)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
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
// 	"net/http" // Still here for HTTP server, but gRPC server needs attention
// 	"time"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection" // <--- IMPORTANT: Added this import!
// 	"movieapp.com/gen"
// 	"movieapp.com/pkg/discovery"
// 	"movieapp.com/pkg/discovery/consul"
// 	"movieapp.com/rating/internal/controller/rating"
// 	grpchandler "movieapp.com/rating/internal/handler/grpc"
// 	// "movieapp.com/rating/internal/ingester/kafka"
// 	"movieapp.com/rating/internal/repository/mysql"
// )

// const serviceName = "rating"

// func main() {
// 	var port int

// 	flag.IntVar(&port, "port", 8082, "API handler port")
// 	flag.Parse()

// 	log.Printf("Starting the rating service on port %d", port)

// 	registry, err := consul.NewRegistry("localhost:8500")
// 	if err != nil {
// 		log.Fatalf("Failed to create registry: %v", err)
// 	}

// 	ctx := context.Background()

// 	instanceID := discovery.GenerateInstanceID(serviceName)

// 	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
// 		log.Fatalf("Failed to register service: %v", err)
// 	}

// 	defer func() {
// 		log.Printf("Deregistering service %s/%s", serviceName, instanceID)
// 		if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
// 			log.Printf("Failed to deregister service: %v", err)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
// 				log.Println("Failed to report healthy state: " + err.Error())
// 			}
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()

// 	repo, err := mysql.New()
// 	if err != nil {
// 		log.Fatalf("Failed to initialize MySQL repository: %v", err) // More specific error message
// 	}

// 	svc := rating.New(repo, nil) // Assuming nil is fine for the ingester part here
// 	h := grpchandler.New(svc)

// 	// --- gRPC Server Setup ---
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) // Changed to ":%d" for all interfaces
// 	if err != nil {
// 		log.Fatalf("Failed to listen for gRPC: %v", err)
// 	}

// 	grpcServer := grpc.NewServer() // Renamed srv to grpcServer for clarity
// 	gen.RegisterRatingServiceServer(grpcServer, h)

// 	// *** ADD THIS LINE FOR GRPC REFLECTION ***
// 	reflection.Register(grpcServer)

// 	log.Printf("gRPC server for %s service serving on %s", serviceName, lis.Addr())

// 	// Start gRPC server in a goroutine so the main function doesn't block,
// 	// allowing for the HTTP server setup to proceed if needed.
// 	go func() {
// 		if err := grpcServer.Serve(lis); err != nil {
// 			log.Fatalf("Failed to serve gRPC: %v", err) // log.Fatalf to ensure main exits on gRPC server failure
// 		}
// 	}()

// 	// --- HTTP Server Setup (if you intend to run both concurrently) ---
// 	// Your existing HTTP server setup here.
// 	// Note: If you don't intend to expose an HTTP API directly from this service,
// 	// you can remove the http.ListenAndServe call.
// 	log.Printf("HTTP server for %s service serving on :%d (if handlers are registered)", serviceName, port)
// 	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
// 		// If both gRPC and HTTP are on the same port, this will conflict.
// 		// If they are on different ports, this would run.
// 		// Given your `net.Listen("tcp", fmt.Sprintf(":%d", port))` for gRPC,
// 		// the HTTP server on the same port will fail to bind.
// 		// You likely want separate ports for gRPC and HTTP, or only one server.
// 		// For now, it will likely panic due to port already in use by gRPC.
// 		log.Fatalf("Failed to serve HTTP: %v", err)
// 	}

// 	// Keep the main goroutine alive if both servers are running in goroutines
// 	// select {} // Uncomment this if you run both in goroutines and need main to wait
// }
