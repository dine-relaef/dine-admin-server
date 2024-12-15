package services

import (
	models "menu-server/src/models"
	postgres "menu-server/src/config/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePlan handles the creation of a new Plan
func CreatePlan(c *gin.Context) {
	var Plan models.Plan

	// Bind JSON to the Plan model
	if err := c.ShouldBindJSON(&Plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save Plan to the postgres
	if err := postgres.DB.Create(&Plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Plan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plan created successfully",
		"Plan": Plan,
	})
}

// GetPlans retrieves all Plans from the postgres
func GetPlans(c *gin.Context) {
	var Plans []models.Plan

	if err := postgres.DB.Find(&Plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "message" : "Plans Found", "Plans": Plans})
}

// GetPlanByID retrieves a specific Plan by ID
func GetPlanByID(c *gin.Context) {
	id := c.Param("id")
	var Plan models.Plan

	if err := postgres.DB.First(&Plan, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Plan": Plan})
}

// UpdatePlan updates a Plan's information by ID
func UpdatePlan(c *gin.Context) {
	id := c.Param("id")
	var Plan models.Plan

	if err := postgres.DB.First(&Plan, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	var input models.Plan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Model(&Plan).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan updated successfully",
		"Plan": Plan,
	})
}

// DeletePlan removes a Plan by ID
func DeletePlan(c *gin.Context) {
	id := c.Param("id")
	if err := postgres.DB.Delete(&models.Plan{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted successfully"})
}
