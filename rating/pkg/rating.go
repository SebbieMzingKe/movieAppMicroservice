package model

// RecordId defines the record id. together with recordtype identifies unique records across all types
type RecordId string

// RecordType defines the record type. together with recordis identifies unique records across all types
type RecordType string

const (
	RecordTypeMovie = RecordType("movie")
)

// UserID defines the user id
type UserID string

// RatingValue defines a value for a rating record
type RatingValue int

// Rating defines an individual rating created by a user for some record
type Rating struct {
	RecordId   string      `json:"recordId"`
	RecordType string      `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}
