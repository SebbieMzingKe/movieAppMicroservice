package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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