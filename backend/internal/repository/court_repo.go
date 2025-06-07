package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"courtopia-reserve/backend/internal/models"
)

type CourtRepository struct {
	collection *mongo.Collection
}

func NewCourtRepository(db *mongo.Database) *CourtRepository {
	return &CourtRepository{
		collection: db.Collection("courts"),
	}
}

func (r *CourtRepository) FindAll(ctx context.Context) ([]*models.Court, error) {
	var courts []*models.Court

	opts := options.Find().SetSort(bson.M{"court_number": 1})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &courts)
	if err != nil {
		return nil, err
	}

	return courts, nil
}

func (r *CourtRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Court, error) {
	var court models.Court

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&court)
	if err != nil {
		return nil, err
	}

	return &court, nil
}

func (r *CourtRepository) FindByCourtNumber(ctx context.Context, courtNumber int) (*models.Court, error) {
	var court models.Court

	filter := bson.M{"court_number": courtNumber}
	err := r.collection.FindOne(ctx, filter).Decode(&court)
	if err != nil {
		return nil, err
	}

	return &court, nil
}

func (r *CourtRepository) FindActiveCourts(ctx context.Context) ([]*models.Court, error) {
	var courts []*models.Court

	filter := bson.M{"is_active": true}
	opts := options.Find().SetSort(bson.M{"court_number": 1})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &courts)
	if err != nil {
		return nil, err
	}

	return courts, nil
}

func (r *CourtRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, isActive bool) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"is_active": isActive}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
