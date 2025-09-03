package movie

import (
	"context"
	"errors"

	metadatamodel "movieapp.com/metadata/pkg/model"
	"movieapp.com/movie/internal/gateway"
	model "movieapp.com/movie/pkg"
	ratingmodel "movieapp.com/rating/pkg"
)

// ErrNotFound is returned when movie metadata is not found
var ErrNotFound = errors.New("movie metadata not found")

// service interfaces
type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordId ratingmodel.RecordId, recordTYpe ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordId ratingmodel.RecordId, recordTYpe ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

// service controller
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	// Add nil checks
	if ratingGateway == nil {
		panic("ratingGateway cannot be nil")
	}
	if metadataGateway == nil {
		panic("metadataGateway cannot be nil")
	}
	
	return &Controller{ratingGateway: ratingGateway, metadataGateway: metadataGateway}
}

// Get returns aggregated rating and movie metadata
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	// Add nil check for safety
	if c.metadataGateway == nil {
		return nil, errors.New("metadataGateway is nil")
	}
	
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	details := &model.MovieDetails{Metadata: *metadata}
	
	// Add nil check for rating gateway
	if c.ratingGateway == nil {
		return details, nil // Return without rating if gateway is nil
	}
	
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordId(id), ratingmodel.RecordTypeMovie)
	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			// Rating not found is OK, just return details without rating
			return details, nil
		}
		// Other errors should be returned
		return nil, err
	}
	
	details.Rating = &rating
	return details, nil
}