package mongo

import (
	"context"
	"time"

	"ap2final_review_service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	reviewsCollection = "reviews"
)

type reviewRepository struct {
	db *mongo.Database
}

func NewReview(db *mongo.Database) ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) (models.Review, error) {
	collection := r.db.Collection(reviewsCollection)

	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now

	result, err := collection.InsertOne(ctx, review)
	if err != nil {
		return models.Review{}, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		review.ID = oid.Hex()
	}

	return *review, nil
}

func (r *reviewRepository) FindByID(ctx context.Context, id string) (models.Review, error) {
	collection := r.db.Collection(reviewsCollection)

	var review models.Review

	filter := bson.M{"_id": id}
	err := collection.FindOne(ctx, filter).Decode(&review)

	if err == mongo.ErrNoDocuments {
		objectID, convErr := primitive.ObjectIDFromHex(id)
		if convErr == nil {
			filter = bson.M{"_id": objectID}
			err = collection.FindOne(ctx, filter).Decode(&review)
		}
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Review{}, models.ErrReviewNotFound
		}
		return models.Review{}, err
	}

	return review, nil
}

func (r *reviewRepository) Find(ctx context.Context, filter models.ReviewFilter) ([]models.Review, error) {
	collection := r.db.Collection(reviewsCollection)

	query := bson.M{}

	if filter.ID != nil {
		objectID, err := primitive.ObjectIDFromHex(*filter.ID)
		if err == nil {
			query["_id"] = objectID
		}
	}

	if len(filter.IDs) > 0 {
		var objectIDs []primitive.ObjectID
		for _, id := range filter.IDs {
			if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
				objectIDs = append(objectIDs, objectID)
			}
		}
		if len(objectIDs) > 0 {
			query["_id"] = bson.M{"$in": objectIDs}
		}
	}

	if filter.UserID != nil {
		query["user_id"] = *filter.UserID
	}

	if filter.MovieID != nil {
		query["movie_id"] = *filter.MovieID
	}

	if filter.Rating != nil {
		query["rating"] = *filter.Rating
	}

	if filter.MinRating != nil {
		query["rating"] = bson.M{"$gte": *filter.MinRating}
	}

	if filter.MaxRating != nil {
		if existing, ok := query["rating"]; ok {
			if ratingQuery, ok := existing.(bson.M); ok {
				ratingQuery["$lte"] = *filter.MaxRating
			}
		} else {
			query["rating"] = bson.M{"$lte": *filter.MaxRating}
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *reviewRepository) Update(ctx context.Context, id string, update models.ReviewUpdateData) (models.Review, error) {
	collection := r.db.Collection(reviewsCollection)

	updateDoc := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	setDoc := updateDoc["$set"].(bson.M)

	if update.Rating != nil {
		setDoc["rating"] = *update.Rating
	}

	if update.Comment != nil {
		setDoc["comment"] = *update.Comment
	}

	if update.IsDeleted != nil {
		setDoc["is_deleted"] = *update.IsDeleted
	}

	var filter bson.M
	var updatedReview models.Review

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	filter = bson.M{"_id": id}
	err := collection.FindOneAndUpdate(ctx, filter, updateDoc, opts).Decode(&updatedReview)

	if err == mongo.ErrNoDocuments {
		if objectID, convErr := primitive.ObjectIDFromHex(id); convErr == nil {
			filter = bson.M{"_id": objectID}
			err = collection.FindOneAndUpdate(ctx, filter, updateDoc, opts).Decode(&updatedReview)
		}
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Review{}, models.ErrReviewNotFound
		}
		return models.Review{}, err
	}

	return updatedReview, nil
}

func (r *reviewRepository) Delete(ctx context.Context, id string) (models.Review, error) {
	return r.Update(ctx, id, models.ReviewUpdateData{
		IsDeleted: models.BoolPtr(true),
	})
}

func (r *reviewRepository) CheckUserReviewExists(ctx context.Context, userID, movieID string) (bool, error) {
	collection := r.db.Collection(reviewsCollection)

	filter := bson.M{
		"user_id":    userID,
		"movie_id":   movieID,
		"is_deleted": bson.M{"$ne": true},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *reviewRepository) GetAverageRating(ctx context.Context, movieID string) (float64, error) {
	collection := r.db.Collection(reviewsCollection)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"movie_id":   movieID,
				"is_deleted": bson.M{"$ne": true},
			},
		},
		{
			"$group": bson.M{
				"_id":            nil,
				"average_rating": bson.M{"$avg": "$rating"},
				"review_count":   bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		AverageRating float64 `bson:"average_rating"`
		ReviewCount   int     `bson:"review_count"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}

		if result.ReviewCount == 0 {
			return 0, nil
		}

		return result.AverageRating, nil
	}

	return 0, nil
}
