package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	metadatamodel "movieapp.com/metadata/pkg/model"
	"movieapp.com/movie/internal/gateway"
	"movieapp.com/pkg/discovery"
	ratingmodel "movieapp.com/rating/pkg"
)

// Gateway 	defines a movie metadata http gateway
type Gateway struct {
	registry discovery.Registry
}

// GetAggregatedRating implements movie.ratingGateway.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordId ratingmodel.RecordId, recordTYpe ratingmodel.RecordType) (float64, error) {
	panic("unimplemented")
}

// PutRating implements movie.ratingGateway.
func (g *Gateway) PutRating(ctx context.Context, recordId ratingmodel.RecordId, recordTYpe ratingmodel.RecordType, rating *ratingmodel.Rating) error {
	panic("unimplemented")
}

// New creates a movie metadata http gateway
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Get gets a movie metadata by movie id
func (g *Gateway) Get(ctx context.Context, id string) (*metadatamodel.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"
	log.Printf("Calling metadadata service. Request: GET %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v *metadatamodel.Metadata

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}