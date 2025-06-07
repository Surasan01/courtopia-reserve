package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"courtopia-reserve/backend/internal/models"
	"courtopia-reserve/backend/internal/repository"
	"courtopia-reserve/backend/pkg/utils"
)

func (h *Handler) CreateBooking(c *gin.Context) {
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*utils.Claims)

	userID, err := primitive.ObjectIDFromHex(userClaims.Subject)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	layout := "15:04"
	startTimeParsed, err := time.Parse(layout, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format, use HH:MM"})
		return
	}

	endTimeParsed, err := time.Parse(layout, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format, use HH:MM"})
		return
	}

	startTime := time.Date(
		bookingDate.Year(),
		bookingDate.Month(),
		bookingDate.Day(),
		startTimeParsed.Hour(),
		startTimeParsed.Minute(),
		0,
		0,
		bookingDate.Location(),
	)

	endTime := time.Date(
		bookingDate.Year(),
		bookingDate.Month(),
		bookingDate.Day(),
		endTimeParsed.Hour(),
		endTimeParsed.Minute(),
		0,
		0,
		bookingDate.Location(),
	)

	if startTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking time must be in the future"})
		return
	}

	if endTime.Before(startTime) || endTime.Equal(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}

	maxDuration := 2 * time.Hour
	if endTime.Sub(startTime) > maxDuration {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking duration cannot exceed 2 hours"})
		return
	}

	court, err := h.courtRepo.FindByCourtNumber(c.Request.Context(), req.CourtNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
		return
	}

	if !court.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Court is not available for booking"})
		return
	}

	isAvailable, err := h.bookingRepo.IsCourtAvailable(c.Request.Context(), req.CourtNumber, bookingDate, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check court availability"})
		return
	}

	if !isAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Court is not available for the selected time"})
		return
	}

	booking := &models.Booking{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		StudentID:        userClaims.StudentID,
		CourtID:          court.ID,
		CourtNumber:      req.CourtNumber,
		BookingDate:      bookingDate,
		StartTime:        startTime,
		EndTime:          endTime,
		Status:           "active",
		NotificationSent: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		UserEmail:        userClaims.Eamil,
	}

	if err := h.bookingRepo.Create(c.Request.Context(), booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	response := models.BookingResponse{
		ID:          booking.ID.Hex(),
		CourtNumber: booking.CourtNumber,
		BookingDate: req.BookingDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      booking.Status,
		CreatedAt:   booking.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetUserBookings(c *gin.Context) {
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*utils.Claims)

	if err := h.bookingRepo.UpdateCompletedBookings(c.Request.Context()); err != nil {
		log.Printf("Error updating completed bookings: %v", err)
	}

	log.Printf("Fetching bookings for student ID: %s", userClaims.StudentID)

	bookings, err := h.bookingRepo.FindByStudentID(c.Request.Context(), userClaims.StudentID)
	if err != nil {
		log.Printf("Error fetching bookings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	var response []models.BookingResponse
	for _, booking := range bookings {
		response = append(response, models.BookingResponse{
			ID:          booking.ID.Hex(),
			CourtNumber: booking.CourtNumber,
			BookingDate: booking.BookingDate.Format("2006-01-02"),
			StartTime:   booking.StartTime.Format("15:04"),
			EndTime:     booking.EndTime.Format("15:04"),
			Status:      booking.Status,
			CreatedAt:   booking.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CancelBooking(c *gin.Context) {
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*utils.Claims)

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.bookingRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	isOwner := booking.StudentID == userClaims.StudentID
	isAdmin := userClaims.Role == "admin"
	if !isOwner && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to cancel this booking"})
		return
	}

	if booking.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is already cancelled"})
		return
	}

	if err := h.bookingRepo.CancelBooking(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking cancelled successfully",
	})
}

func (h *Handler) GetAllBookings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Admin bookings endpoint"})
}

func (h *Handler) CheckAvailability(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Check availability endpoint"})
}

func SendMail(bookingRepo *repository.BookingRepository, userRepo *repository.UserRepository) {
	log.Println("Starting SendMail function...")

	bookings, err := bookingRepo.FindUpcomingBookings(context.Background(), time.Now().Add(1*time.Minute))
	if err != nil {
		log.Printf("Error fetching upcoming bookings: %v", err)
		return
	}

	log.Printf("Found %d upcoming bookings", len(bookings))

	auth := smtp.PlainAuth(
		"",
		"natthawat48.noi@gmail.com",
		"mjxn rpse favy jidl",      
		"smtp.gmail.com",            
	)

	for _, booking := range bookings {
		log.Printf("Processing booking ID: %s", booking.ID.Hex())

		user, err := userRepo.FindByStudentID(context.Background(), booking.StudentID)
		if err != nil {
			log.Printf("Error fetching user email for StudentID %s: %v", booking.StudentID, err)
			continue
		}

		log.Printf("Sending email to: %s", user.Email)

		msg := []byte(fmt.Sprintf(
			"To: %s\r\nSubject: Upcoming Booking Reminder\r\n\r\nDear user,\n\nThis is a reminder for your upcoming booking:\n\nCourt Number: %d\nDate: %s\nTime: %s - %s\n\nThank you for using Courtminton!",
			user.Email,
			booking.CourtNumber,
			booking.BookingDate.Format("2006-01-02"),
			booking.StartTime.Format("15:04"),
			booking.EndTime.Format("15:04"),
		))

		err = smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			"natthawat48.noi@gmail.com",
			[]string{user.Email},
			msg,
		)

		if err != nil {
			log.Printf("Error sending email to %s: %v", user.Email, err)
			continue
		}

		log.Printf("Email sent successfully to: %s", user.Email)

		booking.NotificationSent = true
		if err := bookingRepo.UpdateBooking(context.Background(), booking); err != nil {
			log.Printf("Error updating notification status for booking ID %s: %v", booking.ID.Hex(), err)
		} else {
			log.Printf("Notification status updated for booking ID: %s", booking.ID.Hex())
		}
	}

	log.Println("SendMail function completed.")
}

func (h *Handler) TriggerEmailNotifications(c *gin.Context) {
	SendMail(h.bookingRepo, h.userRepo)

	c.JSON(http.StatusOK, gin.H{"message": "Email notifications triggered"})
}
