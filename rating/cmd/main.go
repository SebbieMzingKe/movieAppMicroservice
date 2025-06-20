package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
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

	f, err := os.Open("base.yaml")

	if err != nil {
		panic(err)
	}

	var cfg config

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	port = cfg.ApiConfig.Port
	log.Printf("Starting the rating service on port %d", port)
	registry, err := consul.NewRegistry("consul-consul-server:8500")

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