package mongo

import (
	"ap2final_review_service/internal/models"
	"context"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) (models.Review, error)
	FindByID(ctx context.Context, id string) (models.Review, error)
	Find(ctx context.Context, filter models.ReviewFilter) ([]models.Review, error)
	Update(ctx context.Context, id string, update models.ReviewUpdateData) (models.Review, error)
	Delete(ctx context.Context, id string) (models.Review, error)
	CheckUserReviewExists(ctx context.Context, userID, movieID string) (bool, error)
	GetAverageRating(ctx context.Context, movieID string) (float64, error)
}
