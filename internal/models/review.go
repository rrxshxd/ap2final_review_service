package models

import (
	"errors"
	"time"
)

type Review struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	MovieID   string    `bson:"movie_id"`
	Rating    int       `bson:"rating"` // 1-5 stars
	Comment   string    `bson:"comment"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	IsDeleted bool      `bson:"is_deleted"`
}

type ReviewFilter struct {
	ID        *string
	IDs       []string
	UserID    *string
	MovieID   *string
	Rating    *int
	MinRating *int
	MaxRating *int
}

type ReviewUpdateData struct {
	Rating    *int
	Comment   *string
	IsDeleted *bool
}

var (
	ErrReviewNotFound      = errors.New("review not found")
	ErrReviewAlreadyExists = errors.New("user has already reviewed this movie")
	ErrInvalidRating       = errors.New("rating must be between 1 and 5")
	ErrEmptyComment        = errors.New("comment cannot be empty")
	ErrInvalidInput        = errors.New("invalid input data")
)

// Helper functions
func (r *Review) Validate() error {
	if r.Rating < 1 || r.Rating > 5 {
		return ErrInvalidRating
	}
	if r.Comment == "" {
		return ErrEmptyComment
	}
	if r.UserID == "" || r.MovieID == "" {
		return ErrInvalidInput
	}
	return nil
}

// Helper for creating pointers
func IntPtr(i int) *int          { return &i }
func StringPtr(s string) *string { return &s }
func BoolPtr(b bool) *bool       { return &b }
