package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"courtopia-reserve/backend/internal/models"
)

func (h *Handler) GetCourts(c *gin.Context) {
	courts, err := h.courtRepo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courts"})
		return
	}

	c.JSON(http.StatusOK, courts)
}

func (h *Handler) GetCourt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	court, err := h.courtRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
		return
	}

	c.JSON(http.StatusOK, court)
}

func (h *Handler) GetAvailableCourts(c *gin.Context) {
	dateStr := c.Query("date")
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")

	if dateStr == "" || startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date, start time and end time are required"})
		return
	}

	bookingDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	layout := "15:04"
	startTimeParsed, err := time.Parse(layout, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format, use HH:MM"})
		return
	}

	endTimeParsed, err := time.Parse(layout, endTimeStr)
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

	availabilities, err := h.bookingRepo.GetAvailableCourts(
		c.Request.Context(),
		bookingDate,
		startTime,
		endTime,
		h.courtRepo,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check court availability"})
		return
	}

	response := models.AvailabilityResponse{
		BookingDate: dateStr,
		StartTime:   startTimeStr,
		EndTime:     endTimeStr,
		Courts:      availabilities,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateCourtStatus(c *gin.Context) {
	// ดึงค่า ID จาก URL
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	var req struct {
		IsActive bool `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.courtRepo.UpdateStatus(c.Request.Context(), id, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update court status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Court status updated successfully"})
}
