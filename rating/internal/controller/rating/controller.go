package rating

import (
	"context"
	"errors"

	"movieapp.com/rating/internal/repository"
	model "movieapp.com/rating/pkg"
)

// ErrNotFound returned when no ratings are found for a record
var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordId model.RecordId, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordId model.RecordId, recordType model.RecordType, rating *model.Rating) error
}

// type Controller struct {
// 	repo ratingRepository
// }

// New creates a rating service controller
// func New(repo ratingRepository) *Controller {
// 	return &Controller{repo}
// }

type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

type Controller struct {
	repo     ratingRepository
	ingester ratingIngester
}

// New creates a rating service
func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{repo, ingester}
}


// writing and getting an aggregated rating
func (c *Controller) GetAggregatedRating(ctx context.Context, recordId model.RecordId, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordId, recordType)
	if err != nil && err == repository.ErrNotFound {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}
	if len(ratings) == 0 {
		return 0, ErrNotFound
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

func (c *Controller) PutRating(ctx context.Context, recordId model.RecordId, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordId, recordType, rating)
}

// StartIngestion start ingestion of the rating events
func (s *Controller) StartIngestion(ctx context.Context) error {
	ch, err := s.ingester.Ingest(ctx)
	if err != nil {
		return err
	}

	for e := range ch {
		if err := s.PutRating(ctx, e.RecordId, e.RecordType, &model.Rating{
			UserID: e.UserID,
			Value:  e.Value,
		}); err != nil {
			return err
		}
	}
	return nil
}
