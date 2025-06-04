package movie

import (
	"context"
	"errors"

	metadatamodel "movieapp.com/metadata/pkg/model"
	"movieapp.com/movie/internal/gateway"
	model "movieapp.com/movie/pkg"
	ratingmodel "movieapp.com/rating/pkg"
)

// ErrNotFound is returrned when movie metadata is not found
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
	ratingGateway ratingGateway
	metadataGateway metadataGateway
}

// New ceretaes a new movie service controller
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway, metadataGateway}
}

// getting movie rating and metadata
// get returns aggregated rating and movie metadata
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error ){
	metadata, err := c.metadataGateway.Get(ctx, id)

	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	}

	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordId(id), ratingmodel.RecordTypeMovie)

	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}
	return details, nil
}