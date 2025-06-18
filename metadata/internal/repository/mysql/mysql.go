package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"movieapp.com/metadata/internal/repository"
	"movieapp.com/metadata/pkg/model"
)

type Repository struct {
	db *sql.DB
}

// new mysql base repository
func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:password@/movieapp")
	if err != nil {
		return nil, err
	}

	return &Repository{db}, nil
}


// retrive metadata by movie id
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string

	row := r.db.QueryRowContext(ctx, "SELECT title, description, director FROM movies WHERE id = ?", id)
	if err := row.Scan(&title, &description, &director); err != nil {

		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &model.Metadata{
		ID: id,
		Title: title,
		Description: description,
		Director: director,
	}, nil
}


// add movie metadata for a given id
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO movies (id, title, description, director) VALUES (?,?, ?, ?)", id, metadata.Title, metadata.Description, metadata.Director)
	return err
}