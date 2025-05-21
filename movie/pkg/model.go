package model

import model "movieapp.com/metadata/pkg"

// moviedetails include movie metadata and its aggregated rating
type MovieDetails struct {
	Rating *float64 `json:"rating,omitEmpty"`
	Metadata model.Metadata `json:"metadata"`
}