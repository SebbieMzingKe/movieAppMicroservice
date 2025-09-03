package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"movieapp.com/gen"
	"movieapp.com/metadata/internal/controller/metadata"
	grpchandler "movieapp.com/metadata/internal/handler/grpc"
	"movieapp.com/metadata/internal/repository/memory"
	"movieapp.com/pkg/discovery"
	"movieapp.com/pkg/discovery/consul"
	"movieapp.com/pkg/tracing"
)

const serviceName = "metadata"

func main() {

	var port int

	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()

	var cfg serviceConfig

	log.Printf("Starting the metadata service on port %d", port)

	tp, err := tracing.NewOtlpGrpcProvider(context.Background(), cfg.Jaeger.URL, serviceName)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

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

	f, err := os.Open("home/seb/Desktop/projects/movvieApp/metadata/configs/base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	repo := memory.New()
	svc := metadata.New(repo)
	h := grpchandler.New(svc)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	srv.Serve(lis)

	reflection.Register(srv)

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
