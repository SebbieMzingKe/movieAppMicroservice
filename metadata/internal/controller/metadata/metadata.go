package metadata

import (
	"context"
	"errors"

	"movieapp.com/metadata/internal/repository"
	model "movieapp.com/metadata/pkg/model"
)

// ErrNotFound is returned when requested record is not found
var ErrNotFound = errors.New("not found")

type MetadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, id string, m *model.Metadata) error
}

// identify metadata service controller
type Controller struct {
	repo MetadataRepository
}

// New creates a metadata service controller
func New(repo MetadataRepository) *Controller {
	return &Controller{repo}
}

// Get returns a movie metadata by id
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (c *Controller) Put(ctx context.Context, m *model.Metadata) error {
	return c.repo.Put(ctx, m.ID, m)
}