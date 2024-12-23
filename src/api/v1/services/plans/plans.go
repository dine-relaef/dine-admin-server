package services_plan

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreatePlan handles the creation of a new Plan
// @Summary Create a new Plan
// @Description Create a new Plan
// @Tags Plans
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param input body models.AddPlanData true "Plan data"
// @Router /api/v1/plans [post]
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
		"Plan":    Plan,
	})
}

// AddPlanFeature adds a feature to a Plan
// @Summary Add a feature to a Plan
// @Description Add a feature to a Plan
// @Tags Plans
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param plan_id query string true "Plan ID"
// @Param feature_id query string true "Feature ID"
// @Router /api/v1/plans/add-feature [put]
func AddPlanFeature(c *gin.Context) {
	// Parse query parameters
	planID := c.Query("plan_id")
	featureID := c.Query("feature_id")

	// Validate inputs
	if planID == "" || featureID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both plan_id and feature_id are required"})
		return
	}

	// Convert strings to UUID
	planUUID, err := uuid.FromString(planID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan_id"})
		return
	}

	featureUUID, err := uuid.FromString(featureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feature_id"})
		return
	}

	var Plan models.Plan
	if err := postgres.DB.First(&Plan, "id = ?", planUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	var Feature models.PlanFeature
	if err := postgres.DB.First(&Feature, "id = ?", featureUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feature not found"})
		return
	}

	// Create a new PlanFeatureAssociation instance
	association := models.PlanFeatureAssociation{
		PlanID:    planUUID,
		FeatureID: featureUUID,
		Plan:      Plan,
		Feature:   Feature,
	}

	if err := postgres.DB.Create(&association).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create association", "details": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "PlanFeatureAssociation added successfully",
		"data":    association,
	})
}

// GelAllPlans retrieves all Plans from the postgres
// @Summary Get all Plans
// @Description Get all Plans
// @Tags Plans
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/v1/plans/all [get]
func GetAllPlans(c *gin.Context) {
	var Plans []models.Plan

	// Eager load the related PlanFeatureAssociations and PlanFeature
	if err := postgres.DB.Preload("PlanFeatureAssociations.Feature").Find(&Plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Plans"})
		return
	}

	if len(Plans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Plans found"})
		return
	}

	var planResponses []models.PlanResponse

	for _, plan := range Plans {
		var planResponse models.PlanResponse
		planResponse.ID = plan.ID
		planResponse.Name = plan.Name
		planResponse.Description = plan.Description
		planResponse.Price = plan.Price
		planResponse.IsActive = plan.IsActive
		planResponse.TrialPeriod = plan.TrialPeriod
		planResponse.CreatedAt = plan.CreatedAt
		planResponse.UpdatedAt = plan.UpdatedAt

		// Iterate over preloaded associations to gather features
		for _, association := range plan.PlanFeatureAssociations {

			planResponse.Feature = append(planResponse.Feature, models.PlanFeature{
				ID:          association.Feature.ID,
				Name:        association.Feature.Name,
				Description: association.Feature.Description,
				CreatedAt:   association.Feature.CreatedAt,
				UpdatedAt:   association.Feature.UpdatedAt,
			})

		}

		planResponses = append(planResponses, planResponse)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plans Found", "Plans": planResponses})
}

// GetPlans retrieves Plans from the postgres
// @Summary Get Plans
// @Description Get Plans
// @Tags Plans
// @Accept json
// @Produce json
// @Router /api/v1/plans [get]
func GetPlans(c *gin.Context) {
	var Plans []models.Plan

	if err := postgres.DB.Preload("PlanFeatureAssociations.Feature").Where("is_active = ?", true).Find(&Plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Plans"})
		return
	}

	if len(Plans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Plans found"})
		return
	}

	var planResponses []models.PlanResponse

	for _, plan := range Plans {
		var planResponse models.PlanResponse
		planResponse.ID = plan.ID
		planResponse.Name = plan.Name
		planResponse.Description = plan.Description
		planResponse.Price = plan.Price
		planResponse.IsActive = plan.IsActive
		planResponse.TrialPeriod = plan.TrialPeriod
		planResponse.CreatedAt = plan.CreatedAt
		planResponse.UpdatedAt = plan.UpdatedAt

		// Iterate over preloaded associations to gather features
		for _, association := range plan.PlanFeatureAssociations {

			planResponse.Feature = append(planResponse.Feature, models.PlanFeature{
				ID:          association.Feature.ID,
				Name:        association.Feature.Name,
				Description: association.Feature.Description,
				CreatedAt:   association.Feature.CreatedAt,
				UpdatedAt:   association.Feature.UpdatedAt,
			})

		}

		planResponses = append(planResponses, planResponse)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plans Found", "Plans": planResponses})
}

// GetPlanByID retrieves a specific Plan by ID
// @Summary Get a Plan by ID
// @Description Get a Plan by ID
// @Tags Plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Router /api/v1/plans/{id} [get]
func GetPlanByID(c *gin.Context) {
	id := c.Param("id")
	var Plan models.Plan

	if err := postgres.DB.Preload("PlanFeatureAssociations.Feature").First(&Plan, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	var planResponse models.PlanResponse
	planResponse.ID = Plan.ID
	planResponse.Name = Plan.Name
	planResponse.Description = Plan.Description
	planResponse.Price = Plan.Price
	planResponse.IsActive = Plan.IsActive
	planResponse.TrialPeriod = Plan.TrialPeriod
	planResponse.CreatedAt = Plan.CreatedAt
	planResponse.UpdatedAt = Plan.UpdatedAt

	// Iterate over preloaded associations to gather features
	for _, association := range Plan.PlanFeatureAssociations {

		planResponse.Feature = append(planResponse.Feature, models.PlanFeature{
			ID:          association.Feature.ID,
			Name:        association.Feature.Name,
			Description: association.Feature.Description,
			CreatedAt:   association.Feature.CreatedAt,
			UpdatedAt:   association.Feature.UpdatedAt,
		})

	}

	c.JSON(http.StatusOK, gin.H{"Plan": planResponse})
}

// UpdatePlan updates a Plan's information by ID
// @Summary Update a Plan
// @Description Update a Plan by ID
// @Tags Plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param Plan body models.UpdatePlanData true "Plan data"
// @Security ApiKeyAuth
// @Router /api/v1/plans/{id} [put]
func UpdatePlan(c *gin.Context) {

	id := c.Param("id")

	var input models.UpdatePlanData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var Plan models.Plan
	if err := postgres.DB.First(&Plan, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	if input.Name != "" {
		Plan.Name = input.Name
	}
	if input.Description != "" {
		Plan.Description = input.Description
	}
	if input.Price != 0 {
		Plan.Price = input.Price
	}
	if input.IsActive != Plan.IsActive {
		Plan.IsActive = input.IsActive
	}
	if input.TrialPeriod != Plan.TrialPeriod {
		Plan.TrialPeriod = input.TrialPeriod
	}

	if err := postgres.DB.Save(&Plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan updated successfully",
		"Plan":    Plan,
	})
}

// DeletePlan deletes a Plan by ID
// @Summary Delete a Plan
// @Description Delete a Plan by ID
// @Tags Plans
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Plan ID"
// @Router /api/v1/plans/{id} [delete]
func DeletePlan(c *gin.Context) {
	id := c.Param("id")
	if err := postgres.DB.Delete(&models.Plan{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted successfully"})
}

// RemovePlanFeature removes a feature from a Plan
// @Summary Remove a feature from a Plan
// @Description Remove a feature from a Plan
// @Tags Plans
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param plan_id query string true "Plan ID"
// @Param feature_id query string true "Feature ID"
// @Router /api/v1/plans/remove-feature [put]
func RemovePlanFeature(c *gin.Context) {
	// Parse query parameters
	planID := c.Query("plan_id")
	featureID := c.Query("feature_id")

	// Validate inputs
	if planID == "" || featureID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both plan_id and feature_id are required"})
		return
	}

	// Convert strings to UUID
	planUUID, err := uuid.FromString(planID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan_id"})
		return
	}

	featureUUID, err := uuid.FromString(featureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feature_id"})
		return
	}

	// Delete the association
	if err := postgres.DB.Where("plan_id = ? AND feature_id = ?", planUUID, featureUUID).Delete(&models.PlanFeatureAssociation{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove feature from Plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feature removed from Plan successfully"})
}
