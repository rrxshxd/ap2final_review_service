package grpc

import (
	"ap2final_review_service/internal/adapter/grpc/dto"
	"context"
	"github.com/sorawaslocked/ap2final_protos_gen/base"
	svc "github.com/sorawaslocked/ap2final_protos_gen/service/review"
	"log/slog"
)

type ReviewServer struct {
	uc  ReviewUseCase
	log *slog.Logger
	svc.UnimplementedReviewServiceServer
}

func NewReviewServer(
	uc ReviewUseCase,
	log *slog.Logger,
) *ReviewServer {
	return &ReviewServer{
		uc:  uc,
		log: log,
	}
}

func (s *ReviewServer) Create(ctx context.Context, req *svc.CreateRequest) (*svc.CreateResponse, error) {
	review := dto.ToReviewFromCreateRequest(req)

	createdReview, err := s.uc.Create(ctx, review)
	if err != nil {
		s.logError("create", err)
		return nil, dto.FromError(err)
	}

	return &svc.CreateResponse{
		Review: dto.FromReviewToPb(createdReview),
	}, nil
}

func (s *ReviewServer) Get(ctx context.Context, req *svc.GetRequest) (*svc.GetResponse, error) {
	review, err := s.uc.GetByID(ctx, req.ID)
	if err != nil {
		s.logError("get", err)
		return nil, dto.FromError(err)
	}

	return &svc.GetResponse{
		Review: dto.FromReviewToPb(review),
	}, nil
}

func (s *ReviewServer) GetAll(ctx context.Context, req *svc.GetAllRequest) (*svc.GetAllResponse, error) {
	reviews, err := s.uc.GetAll(ctx)
	if err != nil {
		s.logError("get all", err)
		return nil, dto.FromError(err)
	}

	var reviewsPb []*base.Review
	for _, review := range reviews {
		reviewsPb = append(reviewsPb, dto.FromReviewToPb(review))
	}

	return &svc.GetAllResponse{
		Reviews: reviewsPb,
	}, nil
}

func (s *ReviewServer) GetByUser(ctx context.Context, req *svc.GetByUserRequest) (*svc.GetByUserResponse, error) {
	reviews, err := s.uc.GetByUserID(ctx, req.UserID)
	if err != nil {
		s.logError("get by user", err)
		return nil, dto.FromError(err)
	}

	var reviewsPb []*base.Review
	for _, review := range reviews {
		reviewsPb = append(reviewsPb, dto.FromReviewToPb(review))
	}

	return &svc.GetByUserResponse{
		Reviews: reviewsPb,
	}, nil
}

func (s *ReviewServer) GetByMovie(ctx context.Context, req *svc.GetByMovieRequest) (*svc.GetByMovieResponse, error) {
	reviews, err := s.uc.GetByMovieID(ctx, req.MovieID)
	if err != nil {
		s.logError("get by movie", err)
		return nil, dto.FromError(err)
	}

	var reviewsPb []*base.Review
	for _, review := range reviews {
		reviewsPb = append(reviewsPb, dto.FromReviewToPb(review))
	}

	return &svc.GetByMovieResponse{
		Reviews: reviewsPb,
	}, nil
}

func (s *ReviewServer) Update(ctx context.Context, req *svc.UpdateRequest) (*svc.UpdateResponse, error) {
	id, update := dto.ToReviewUpdateFromUpdateRequest(req)

	updatedReview, err := s.uc.UpdateByID(ctx, id, update)
	if err != nil {
		s.logError("update", err)
		return nil, dto.FromError(err)
	}

	return &svc.UpdateResponse{
		Review: dto.FromReviewToPb(updatedReview),
	}, nil
}

func (s *ReviewServer) Delete(ctx context.Context, req *svc.DeleteRequest) (*svc.DeleteResponse, error) {
	deletedReview, err := s.uc.DeleteByID(ctx, req.ID)
	if err != nil {
		s.logError("delete", err)
		return nil, dto.FromError(err)
	}

	return &svc.DeleteResponse{
		Review: dto.FromReviewToPb(deletedReview),
	}, nil
}

func (s *ReviewServer) logError(op string, err error) {
	s.log.Error("review operation failed", slog.String("operation", op), slog.String("error", err.Error()))
}
