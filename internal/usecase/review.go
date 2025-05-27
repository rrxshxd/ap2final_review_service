package usecase

import (
	"context"
	"log/slog"
	"time"

	"ap2final_review_service/internal/models"
)

type reviewUseCase struct {
	repo ReviewRepository
	log  *slog.Logger
}

func NewReviewUseCase(repo ReviewRepository, log *slog.Logger) ReviewUseCase {
	return &reviewUseCase{
		repo: repo,
		log:  log,
	}
}

func (uc *reviewUseCase) Create(ctx context.Context, review models.Review) (models.Review, error) {
	if err := review.Validate(); err != nil {
		return models.Review{}, err
	}

	exists, err := uc.repo.CheckUserReviewExists(ctx, review.UserID, review.MovieID)
	if err != nil {
		return models.Review{}, err
	}
	if exists {
		return models.Review{}, models.ErrReviewAlreadyExists
	}

	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now
	review.IsDeleted = false

	createdReview, err := uc.repo.Create(ctx, &review)
	if err != nil {
		uc.log.Error("failed to create review", "error", err)
		return models.Review{}, err
	}

	return createdReview, nil
}

func (uc *reviewUseCase) GetByID(ctx context.Context, id string) (models.Review, error) {
	review, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return models.Review{}, err
	}

	if review.IsDeleted {
		return models.Review{}, models.ErrReviewNotFound
	}

	return review, nil
}

func (uc *reviewUseCase) GetAll(ctx context.Context) ([]models.Review, error) {
	reviews, err := uc.repo.Find(ctx, models.ReviewFilter{})
	if err != nil {
		return nil, err
	}

	var activeReviews []models.Review
	for _, review := range reviews {
		if !review.IsDeleted {
			activeReviews = append(activeReviews, review)
		}
	}

	return activeReviews, nil
}

func (uc *reviewUseCase) GetByUserID(ctx context.Context, userID string) ([]models.Review, error) {
	reviews, err := uc.repo.Find(ctx, models.ReviewFilter{UserID: &userID})
	if err != nil {
		return nil, err
	}

	var activeReviews []models.Review
	for _, review := range reviews {
		if !review.IsDeleted {
			activeReviews = append(activeReviews, review)
		}
	}

	return activeReviews, nil
}

func (uc *reviewUseCase) GetByMovieID(ctx context.Context, movieID string) ([]models.Review, error) {
	reviews, err := uc.repo.Find(ctx, models.ReviewFilter{MovieID: &movieID})
	if err != nil {
		return nil, err
	}

	var activeReviews []models.Review
	for _, review := range reviews {
		if !review.IsDeleted {
			activeReviews = append(activeReviews, review)
		}
	}

	return activeReviews, nil
}

func (uc *reviewUseCase) UpdateByID(ctx context.Context, id string, update models.ReviewUpdateData) (models.Review, error) {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return models.Review{}, err
	}

	if existing.IsDeleted {
		return models.Review{}, models.ErrReviewNotFound
	}

	if update.Rating != nil {
		if *update.Rating < 1 || *update.Rating > 5 {
			return models.Review{}, models.ErrInvalidRating
		}
	}

	if update.Comment != nil && *update.Comment == "" {
		return models.Review{}, models.ErrEmptyComment
	}

	updatedReview, err := uc.repo.Update(ctx, id, update)
	if err != nil {
		uc.log.Error("failed to update review", "review_id", id, "error", err)
		return models.Review{}, err
	}

	return updatedReview, nil
}

func (uc *reviewUseCase) DeleteByID(ctx context.Context, id string) (models.Review, error) {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return models.Review{}, err
	}

	if existing.IsDeleted {
		return models.Review{}, models.ErrReviewNotFound
	}

	update := models.ReviewUpdateData{
		IsDeleted: models.BoolPtr(true),
	}

	deletedReview, err := uc.repo.Update(ctx, id, update)
	if err != nil {
		uc.log.Error("failed to delete review", "review_id", id, "error", err)
		return models.Review{}, err
	}

	return deletedReview, nil
}

func (uc *reviewUseCase) GetMovieAverageRating(ctx context.Context, movieID string) (float64, error) {
	average, err := uc.repo.GetAverageRating(ctx, movieID)
	if err != nil {
		return 0, err
	}

	return average, nil
}
