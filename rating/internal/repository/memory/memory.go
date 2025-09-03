package memory

import (
	"context"

	"go.opentelemetry.io/otel"
	"movieapp.com/rating/internal/repository"
	model "movieapp.com/rating/pkg"
)

const tracerID = "metadata-repository-memory"

type Repository struct {
	data map[model.RecordType]map[model.RecordId][]model.Rating
}

func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordId][]model.Rating{}}
}

// Get retrives all ratings for a given record
func (r *Repository) Get(ctx context.Context, recordId model.RecordId, recordType model.RecordType) ([]model.Rating, error) {
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Get")
	defer span.End()
	typeMap, ok := r.data[recordType]
	if !ok {
		return nil, repository.ErrNotFound
	}

	ratings, ok := typeMap[recordId]

	if !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	// if _, ok := r.data[recordType]; !ok {
	// 	return nil, repository.ErrNotFound
	// }

	// if ratings, ok := r.data[recordType][recordId]; ok || len(ratings) == 0 {
	// 	return nil, repository.ErrNotFound
	// }

	return ratings, nil
}

// Put adds a rating for a given record
func (r *Repository) Put(ctx context.Context, recordId model.RecordId, recordType model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordId][]model.Rating{}

	}

	r.data[recordType][recordId] = append(r.data[recordType][recordId], *rating)
	return nil
}
