package usecase

import (
	"ap2final_review_service/internal/models"
	"context"
)

type ReviewUseCase interface {
	Create(ctx context.Context, review models.Review) (models.Review, error)
	GetByID(ctx context.Context, id string) (models.Review, error)
	GetAll(ctx context.Context) ([]models.Review, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Review, error)
	GetByMovieID(ctx context.Context, movieID string) ([]models.Review, error)
	UpdateByID(ctx context.Context, id string, update models.ReviewUpdateData) (models.Review, error)
	DeleteByID(ctx context.Context, id string) (models.Review, error)
	GetMovieAverageRating(ctx context.Context, movieID string) (float64, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) (models.Review, error)
	FindByID(ctx context.Context, id string) (models.Review, error)
	Find(ctx context.Context, filter models.ReviewFilter) ([]models.Review, error)
	Update(ctx context.Context, id string, update models.ReviewUpdateData) (models.Review, error)
	Delete(ctx context.Context, id string) (models.Review, error)
	CheckUserReviewExists(ctx context.Context, userID, movieID string) (bool, error)
	GetAverageRating(ctx context.Context, movieID string) (float64, error)
}
