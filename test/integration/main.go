package main

import (
	"context"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"movieapp.com/gen"
	"movieapp.com/pkg/discovery"
	memory "movieapp.com/pkg/discovery/memorypackage"
)

const (
	metadataSvcName = "metadata"
	ratingSvcName   = "rating"
	movieSvcName    = "movie"

	metadataSvcAddr = "localhost:8081"
	ratingSvcAddr   = "localhost:8082"
	movieSvcAddr    = "localhost:8083"
)

func main() {
	log.Println("Starting the integration tests")
	ctx := context.Background()
	registry := memory.NewRegistry()
	log.Println("seting up service handlers and clients")

	metadataSrv := startMetadataService(ctx, registry)
	defer metadataSrv.GracefulStop()

	ratingSrv := startRatingService(ctx, registry)
	defer ratingSrv.GracefulStop()

	movieSrv := startMovieService(ctx, registry)
	defer movieSrv.GracefulStop()

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	metadataConn, err := grpc.Dial(metadataSvcAddr, opts)

	if err != nil {
		panic(err)
	}

	defer metadataConn.Close()

	metadataClient := gen.NewMetadataServiceClient(metadataConn)

	ratingConn, err := grpc.Dial(ratingSvcAddr, opts)

	if err != nil {
		panic(err)
	}

	defer ratingConn.Close()

	ratingClient := gen.NewRatingServiceClient(ratingConn)

	movieConn, err := grpc.Dial(ratingSvcAddr, opts)

	if err != nil {
		panic(err)
	}

	defer movieConn.Close()

	movieClient := gen.NewMovieServiceClient(movieConn)

	log.Println("saving test metadata via metadata service")

	m := &gen.Metadata{
		Id:          "the movie",
		Title:       "the movie",
		Description: "the movie the one and only",
		Director:    "Mr. Seb",
	}

	if _, err := metadataClient.PutMetadata(ctx, &gen.PutMetadataRequest{
		Metadata: m,
	}); err != nil {
		log.Fatalf("put metadata: %v", err)
	}

	log.Println("retrieving test metadata via metadata service")

	getMetadataResp, err := metadataClient.GetMetadata(ctx, &gen.GetMetadataRequest{
		MovieId: m.Id,
	})

	if err != nil {
		log.Fatalf("get metdata: %v", err)
	}

	if diff := cmp.Diff(getMetadataResp.Metadata, m, cmpopts.IgnoreUnexported(gen.Metadata{})); diff != "" {
		log.Fatalf("get metadata after put mismatch: %v", diff)
	}

	log.Println("Getting movie details via movie service")

	wantMovieDetails := &gen.MovieDetails{
		Metadata: m,
	}

	getMovieDetailResp, err := movieClient.GetMovieDetails(ctx, &gen.GetMovieDetailsRequest{MovieId: m.Id})

	if err != nil {
		log.Fatalf("get movie details: %v", err)
	}

	if diff := cmp.Diff(getMovieDetailResp.MovieDetails, wantMovieDetails, cmpopts.IgnoreUnexported(gen.MovieDetails{}, gen.Metadata{})); diff != "" {
		log.Fatalf("get movie details after put mismatch: %v", err)
	}

	log.Println("Saving first rating via rating service")

	const userID = "user0"
	const recordTypeMovie = "movie"

	firstRating := int32(5)

	if _, err = ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:      userID,
		RecordId:    m.Id,
		RecordType:  recordTypeMovie,
		RatingValue: firstRating,
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}

	log.Println("retrieving initial aggregated rating via rating srvice")
	getAggregatedRatingResp, err := ratingClient.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   m.Id,
		RecordType: recordTypeMovie,
	})

	if err != nil {
		log.Fatalf("get aggregated rating: %v", err)
	}

	if got, want := getAggregatedRatingResp.RatingValue, float64(5); got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}

	log.Println("Saving second rating via rating service")

	secondRating := int32(1)

	if _, err = ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:      userID,
		RecordId:    m.Id,
		RecordType:  recordTypeMovie,
		RatingValue: firstRating,
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}

	log.Println("saving new aggregated rating via rating srvice")
	getAggregatedRatingResp, err = ratingClient.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   m.Id,
		RecordType: recordTypeMovie,
	})

	if err != nil {
		log.Fatalf("get aggregated rating: %v", err)
	}

	wantRating := float64((firstRating + secondRating) / 2)
	if got, want := getAggregatedRatingResp.RatingValue, wantRating; got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}

	log.Println("getting updated movie details via movie service")
	getMovieDetailResp, err = movieClient.GetMovieDetails(ctx, &gen.GetMovieDetailsRequest{MovieId: m.Id})

	if err != nil {
		log.Fatalf("get movie details: %v", err)
	}

	wantMovieDetails.Rating = float32(wantRating)

	if diff := cmp.Diff(getMovieDetailResp.MovieDetails, wantMovieDetails, cmpopts.IgnoreUnexported(gen.MovieDetails{}, gen.Metadata{})); diff != "" {
		log.Fatalf("get movie dtails after update mismatch: %v", err)
	}
	
	log.Println("integration test execution successsful")

}

func startMetadataService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	return nil
}

func startRatingService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	return nil
}

func startMovieService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	return nil
}
