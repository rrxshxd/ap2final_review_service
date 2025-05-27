package dto

import (
	"ap2final_review_service/internal/models"
	"github.com/sorawaslocked/ap2final_protos_gen/base"
	svc "github.com/sorawaslocked/ap2final_protos_gen/service/review"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToReviewFromCreateRequest(req *svc.CreateRequest) models.Review {
	return models.Review{
		UserID:  req.UserID,
		MovieID: req.MovieID,
		Rating:  int(req.Rating),
		Comment: req.Comment,
	}
}

func ToReviewUpdateFromUpdateRequest(req *svc.UpdateRequest) (string, models.ReviewUpdateData) {
	update := models.ReviewUpdateData{}

	if req.Rating != nil {
		rating := int(*req.Rating)
		update.Rating = &rating
	}

	if req.Comment != nil {
		update.Comment = req.Comment
	}

	if req.IsDeleted != nil {
		update.IsDeleted = req.IsDeleted
	}

	return req.ID, update
}

func FromReviewToPb(review models.Review) *base.Review {
	return &base.Review{
		ID:        review.ID,
		UserID:    review.UserID,
		MovieID:   review.MovieID,
		Rating:    int32(review.Rating),
		Comment:   review.Comment,
		CreatedAt: timestamppb.New(review.CreatedAt),
		UpdatedAt: timestamppb.New(review.UpdatedAt),
		IsDeleted: review.IsDeleted,
	}
}
