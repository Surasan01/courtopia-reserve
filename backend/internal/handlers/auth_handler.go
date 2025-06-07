package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"courtopia-reserve/backend/internal/models"
	"courtopia-reserve/backend/pkg/utils"
)

func (h *Handler) Register(c *gin.Context) {

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := h.userRepo.FindByStudentID(c.Request.Context(), req.StudentID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "รหัสนักศึกษานี้ถูกใช้งานแล้ว"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	user := &models.User{
		ID:        primitive.NewObjectID(),
		StudentID: req.StudentID,
		Password:  hashedPassword,
		Name:      req.Name,
		Email:     req.Email,
		Role:      "user", 
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userRepo.FindByStudentID(c.Request.Context(), req.StudentID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสนักศึกษาหรือรหัสผ่านไม่ถูกต้อง"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสนักศึกษาหรือรหัสผ่านไม่ถูกต้อง"})
		return
	}

	token, err := utils.GenerateToken(user, h.jwtSecret, 24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:     token,
		StudentID: user.StudentID,
		Name:      user.Name,
		Role:      user.Role,
	})
}
