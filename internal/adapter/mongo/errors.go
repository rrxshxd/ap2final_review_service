package mongo

import (
	"ap2final_review_service/internal/models"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleMongoError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.ErrReviewNotFound
	}

	if mongo.IsDuplicateKeyError(err) {
		return models.ErrReviewAlreadyExists
	}

	return err
}

func IsDuplicateError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}
