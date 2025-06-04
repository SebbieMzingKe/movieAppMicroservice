package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	metadataModel "movieapp.com/metadata/pkg/model"
	"movieapp.com/movie/internal/gateway"
	"movieapp.com/pkg/discovery"
	ratingModel "movieapp.com/rating/pkg"
)

type Gateway struct {
	registry discovery.Registry
}

// Get implements movie.metadataGateway.
func (g *Gateway) Get(ctx context.Context, id string) (*metadataModel.Metadata, error) {
	panic("unimplemented")
}

// New creates a new http gateway for a rating service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordId ratingModel.RecordId, recordType ratingModel.RecordType) (float64, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	// req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)

	if err != nil {
		return 0, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	log.Printf("Calling rating service. Request: GET %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v float64

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// Rating creation request
// PutRating writes a rating
func (g *Gateway) PutRating(ctx context.Context, recordId ratingModel.RecordId, recordType ratingModel.RecordType, rating *ratingModel.Rating) error {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")

	if err != nil {
		return err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	log.Printf("Calling the rating service. Request: PUT %s", url)
	req, err := http.NewRequest(http.MethodPut, addrs[0], nil)

	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
