package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StudentID string             `bson:"student_id" json:"studentId"` // ใช้เป็น username ในการ login
	Password  string             `bson:"password" json:"-"`           // ไม่ส่ง password กลับไป
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty"` // optional
	Role      string             `bson:"role" json:"role"`                       // user, admin
	ProfilePicture string        `bson:"profile_picture,omitempty" json:"profilePicture,omitempty"` // URL ของรูปโปรไฟล์
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Court represents a badminton court
type Court struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CourtNumber int                `bson:"court_number" json:"courtNumber"` // เลขคอร์ท 1-6
	Name        string             `bson:"name" json:"name"`
	IsActive    bool               `bson:"is_active" json:"isActive"`                    // สถานะว่าใช้งานได้หรือไม่
	Location    string             `bson:"location,omitempty" json:"location,omitempty"` // optional
}

// Booking represents a court booking
type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id" json:"userId"`
	StudentID   string             `bson:"student_id" json:"studentId"` // เก็บ StudentID ไว้ด้วยเพื่อง่ายต่อการค้นหา
	CourtID     primitive.ObjectID `bson:"court_id" json:"courtId"`
	CourtNumber int                `bson:"court_number" json:"courtNumber"` // เก็บเลขคอร์ทไว้ด้วยเพื่อความสะดวก
	BookingDate time.Time          `bson:"booking_date" json:"bookingDate"` // วันที่จอง
	StartTime   time.Time          `bson:"start_time" json:"startTime"`     // เวลาเริ่มใช้คอร์ท
	EndTime     time.Time          `bson:"end_time" json:"endTime"`         // เวลาสิ้นสุด (ไม่เกิน 2 ชั่วโมงจากเวลาเริ่ม)
	Status      string             `bson:"status" json:"status"`            // active, cancelled
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	NotificationSent bool `bson:"notification_sent"`
	UserEmail        string             `bson:"user_email" json:"userEmail"`
}

// DTO objects (Data Transfer Objects) for requests and responses

// RegisterRequest represents the data needed for user registration
type RegisterRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email,omitempty"`
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// LoginResponse represents the data returned after successful login
type LoginResponse struct {
	Token     string `json:"token"`
	StudentID string `json:"studentId"`
	Name      string `json:"name"`
	Role      string `json:"role"`
}

// BookingRequest represents the data needed to create a booking
type BookingRequest struct {
	CourtNumber int    `json:"courtNumber" binding:"required"`
	BookingDate string `json:"bookingDate" binding:"required"` // Format: YYYY-MM-DD
	StartTime   string `json:"startTime" binding:"required"`   // Format: HH:MM
	EndTime     string `json:"endTime" binding:"required"`     // Format: HH:MM
}

// BookingResponse represents a booking with additional information
type BookingResponse struct {
	ID          string    `json:"id"`
	CourtNumber int       `json:"courtNumber"`
	BookingDate string    `json:"bookingDate"` // Format: YYYY-MM-DD
	StartTime   string    `json:"startTime"`   // Format: HH:MM
	EndTime     string    `json:"endTime"`     // Format: HH:MM
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// AvailabilityRequest represents the data needed to check court availability
type AvailabilityRequest struct {
	CourtNumber int    `json:"courtNumber,omitempty"`          // Optional, all courts if not provided
	BookingDate string `json:"bookingDate" binding:"required"` // Format: YYYY-MM-DD
	StartTime   string `json:"startTime" binding:"required"`   // Format: HH:MM
	EndTime     string `json:"endTime" binding:"required"`     // Format: HH:MM
}

// CourtAvailability represents the availability of a specific court
type CourtAvailability struct {
	CourtNumber int  `json:"courtNumber"`
	IsAvailable bool `json:"isAvailable"`
}

// AvailabilityResponse represents all available courts for a specific time
type AvailabilityResponse struct {
	BookingDate string               `json:"bookingDate"`
	StartTime   string               `json:"startTime"`
	EndTime     string               `json:"endTime"`
	Courts      []*CourtAvailability `json:"courts"` // เปลี่ยนจาก []CourtAvailability เป็น []*CourtAvailability
}
