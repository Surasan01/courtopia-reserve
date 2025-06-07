package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StudentID string             `bson:"student_id" json:"studentId"` 
	Password  string             `bson:"password" json:"-"`         
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty"` 
	Role      string             `bson:"role" json:"role"`                       
	ProfilePicture string        `bson:"profile_picture,omitempty" json:"profilePicture,omitempty"` 
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

type Court struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CourtNumber int                `bson:"court_number" json:"courtNumber"` 
	Name        string             `bson:"name" json:"name"`
	IsActive    bool               `bson:"is_active" json:"isActive"`                    
	Location    string             `bson:"location,omitempty" json:"location,omitempty"` 
}

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id" json:"userId"`
	StudentID   string             `bson:"student_id" json:"studentId"`
	CourtID     primitive.ObjectID `bson:"court_id" json:"courtId"`
	CourtNumber int                `bson:"court_number" json:"courtNumber"`
	BookingDate time.Time          `bson:"booking_date" json:"bookingDate"` 
	StartTime   time.Time          `bson:"start_time" json:"startTime"`     
	EndTime     time.Time          `bson:"end_time" json:"endTime"`       
	Status      string             `bson:"status" json:"status"`          
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	NotificationSent bool `bson:"notification_sent"`
	UserEmail        string             `bson:"user_email" json:"userEmail"`
}

type RegisterRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email,omitempty"`
}

type LoginRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	StudentID string `json:"studentId"`
	Name      string `json:"name"`
	Role      string `json:"role"`
}

type BookingRequest struct {
	CourtNumber int    `json:"courtNumber" binding:"required"`
	BookingDate string `json:"bookingDate" binding:"required"` 
	StartTime   string `json:"startTime" binding:"required"`   
	EndTime     string `json:"endTime" binding:"required"`     
}

type BookingResponse struct {
	ID          string    `json:"id"`
	CourtNumber int       `json:"courtNumber"`
	BookingDate string    `json:"bookingDate"` 
	StartTime   string    `json:"startTime"`   
	EndTime     string    `json:"endTime"`     
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AvailabilityRequest struct {
	CourtNumber int    `json:"courtNumber,omitempty"`         
	BookingDate string `json:"bookingDate" binding:"required"` 
	StartTime   string `json:"startTime" binding:"required"`   
	EndTime     string `json:"endTime" binding:"required"`    
}

type CourtAvailability struct {
	CourtNumber int  `json:"courtNumber"`
	IsAvailable bool `json:"isAvailable"`
}

type AvailabilityResponse struct {
	BookingDate string               `json:"bookingDate"`
	StartTime   string               `json:"startTime"`
	EndTime     string               `json:"endTime"`
	Courts      []*CourtAvailability `json:"courts"` 
}
