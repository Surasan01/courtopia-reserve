package handlers

import (
	"fmt"
	"net/http"
	"time"

	"courtopia-reserve/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) GetProfile(c *gin.Context) {
	claims := c.MustGet("user").(*utils.Claims)

	user, err := h.userRepo.FindByStudentID(c.Request.Context(), claims.StudentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"studentId":      user.StudentID,
		"name":           user.Name,
		"email":          user.Email,
		"role":           user.Role,
		"profilePicture": user.ProfilePicture,
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	claims := c.MustGet("user").(*utils.Claims)

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	filter := bson.M{"student_id": claims.StudentID}
	update := bson.M{
		"$set": bson.M{
			"name":       req.Name,
			"email":      req.Email,
			"updated_at": time.Now(),
		},
	}

	err := h.userRepo.UpdateOne(c.Request.Context(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *Handler) UploadProfilePicture(c *gin.Context) {
	claims := c.MustGet("user").(*utils.Claims)

	file, err := c.FormFile("profilePicture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	filename := fmt.Sprintf("%s_%d_%s", claims.StudentID, time.Now().Unix(), file.Filename)

	filePath := fmt.Sprintf("uploads/profile_pictures/%s", filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	baseURL := "http://localhost:8000" 
	fullURL := fmt.Sprintf("%s/%s", baseURL, filePath)

	filter := bson.M{"student_id": claims.StudentID}
	update := bson.M{
		"$set": bson.M{
			"profile_picture": fullURL,
			"updated_at":      time.Now(),
		},
	}

	if err := h.userRepo.UpdateOne(c.Request.Context(), filter, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile picture uploaded successfully", "profilePicture": fullURL})
}
