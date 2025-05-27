package dto

import (
	"ap2final_review_service/internal/models"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromError(err error) error {
	if errors.Is(err, models.ErrReviewNotFound) {
		return status.Error(codes.NotFound, "review not found")
	}

	if errors.Is(err, models.ErrReviewAlreadyExists) {
		return status.Error(codes.AlreadyExists, "user has already reviewed this movie")
	}

	if errors.Is(err, models.ErrInvalidRating) {
		return status.Error(codes.InvalidArgument, "rating must be between 1 and 5")
	}

	if errors.Is(err, models.ErrEmptyComment) {
		return status.Error(codes.InvalidArgument, "comment cannot be empty")
	}

	if errors.Is(err, models.ErrInvalidInput) {
		return status.Error(codes.InvalidArgument, "invalid input data")
	}

	return status.Error(codes.Internal, "internal server error")
}
