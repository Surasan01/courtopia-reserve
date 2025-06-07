package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"courtopia-reserve/backend/internal/models"
)

type BookingRepository struct {
	collection *mongo.Collection
}

func NewBookingRepository(db *mongo.Database) *BookingRepository {
	return &BookingRepository{
		collection: db.Collection("bookings"),
	}
}

func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	booking.Status = "active"

	_, err := r.collection.InsertOne(ctx, booking)
	return err
}

func (r *BookingRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Booking, error) {
	var booking models.Booking

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&booking)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) FindByStudentID(ctx context.Context, studentID string) ([]*models.Booking, error) {
	var bookings []*models.Booking

	filter := bson.M{"student_id": studentID}

	opts := options.Find().SetSort(bson.D{
		{Key: "booking_date", Value: -1},
		{Key: "start_time", Value: -1},
	})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("MongoDB Find error: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &bookings)
	if err != nil {
		log.Printf("MongoDB cursor.All error: %v", err)
		return nil, err
	}

	if bookings == nil {
		return []*models.Booking{}, nil
	}

	return bookings, nil
}

func (r *BookingRepository) FindActiveBookingsByStudentID(ctx context.Context, studentID string) ([]*models.Booking, error) {
	var bookings []*models.Booking

	filter := bson.M{
		"student_id": studentID,
		"status":     "active",
		"end_time":   bson.M{"$gte": time.Now()},
	}

	opts := options.Find().SetSort(bson.D{
		{Key: "booking_date", Value: 1},
		{Key: "start_time", Value: 1},
	})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &bookings)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *BookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	booking.UpdatedAt = time.Now()

	filter := bson.M{"_id": booking.ID}
	update := bson.M{"$set": booking}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *BookingRepository) CancelBooking(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"status":     "cancelled",
		"updated_at": time.Now(),
	}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *BookingRepository) IsCourtAvailable(ctx context.Context, courtNumber int, bookingDate time.Time, startTime time.Time, endTime time.Time) (bool, error) {
	startOfDay := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), 0, 0, 0, 0, bookingDate.Location())
	endOfDay := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), 23, 59, 59, 999999999, bookingDate.Location())

	filter := bson.M{
		"court_number": courtNumber,
		"booking_date": bson.M{
			"$gte": startOfDay,
			"$lte": endOfDay,
		},
		"status": "active",
		"$or": []bson.M{
			{
				"start_time": bson.M{"$lte": startTime},
				"end_time":   bson.M{"$gt": startTime},
			},
			{
				"start_time": bson.M{"$lt": endTime},
				"end_time":   bson.M{"$gte": endTime},
			},
			{
				"start_time": bson.M{"$gte": startTime},
				"end_time":   bson.M{"$lte": endTime},
			},
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *BookingRepository) GetAvailableCourts(ctx context.Context, bookingDate time.Time, startTime time.Time, endTime time.Time, courtRepo *CourtRepository) ([]*models.CourtAvailability, error) {
	courts, err := courtRepo.FindActiveCourts(ctx)
	if err != nil {
		return nil, err
	}

	var availabilities []*models.CourtAvailability

	for _, court := range courts {
		isAvailable, err := r.IsCourtAvailable(ctx, court.CourtNumber, bookingDate, startTime, endTime)
		if err != nil {
			return nil, err
		}

		availabilities = append(availabilities, &models.CourtAvailability{
			CourtNumber: court.CourtNumber,
			IsAvailable: isAvailable,
		})
	}

	return availabilities, nil
}

func (r *BookingRepository) UpdateCompletedBookings(ctx context.Context) error {
	now := time.Now()

	filter := bson.M{
		"status":   "active",
		"end_time": bson.M{"$lt": now},
	}

	update := bson.M{
		"$set": bson.M{
			"status":     "completed",
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *BookingRepository) FindUpcomingBookings(ctx context.Context, beforeTime time.Time) ([]*models.Booking, error) {
    beforeTime = beforeTime.Truncate(time.Minute)

    filter := bson.M{
        "$expr": bson.M{
            "$and": []bson.M{
                {"$lte": []interface{}{
                    bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d %H:%M", "date": "$start_time"}},
                    beforeTime.Format("2006-01-02 15:04"),
                }},
				{"$eq": []interface{}{"$notification_sent", false}},
            },
        },
    }

    log.Printf("FindUpcomingBookings filter: %+v", filter) 

    var bookings []*models.Booking
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        log.Printf("Error fetching upcoming bookings: %v", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    err = cursor.All(ctx, &bookings)
    if err != nil {
        log.Printf("Error decoding bookings: %v", err)
        return nil, err
    }

    log.Printf("Upcoming bookings found: %+v", bookings) 

    return bookings, nil
}

func (r *BookingRepository) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	filter := bson.M{"_id": booking.ID}
	update := bson.M{"$set": bson.M{"notification_sent": booking.NotificationSent}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating booking notification status: %v", err)
		return err
	}

	return nil
}
