package metadata

import (
	"errors"
	"testing"

	"movieapp.com/metadata/internal/repository"
	"movieapp.com/metadata/pkg/model"
)

func TestMetadata(t *testing.T) {
	tests := []struct {
		name       string
		expRepoRes *model.Metadata
		expRepErr  error
		wantRes    *model.Metadata
		wantErr    error
	}{
		{
			name:      "not found",
			expRepErr: repository.ErrNotFound,
			wantErr:   ErrNotFound,
		},
		{
			name:      "unexpected error",
			expRepErr: errors.New("unexpected error"),
			wantErr:   errors.New("unexpected error"),
		},
		{
			name:       "success",
			expRepoRes: &model.Metadata{},
			wantRes:    &model.Metadata{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
	}
}
