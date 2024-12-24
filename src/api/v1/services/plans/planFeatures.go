package services_plan

import (
	postgres "menu-server/src/config/database"
	models_plan "menu-server/src/models/plans"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePlanFeatures handles the creation of a new PlanFeatures
// @Summary Create a new PlanFeatures
// @Description Create a new PlanFeatures
// @Tags PlanFeatures
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param input body models_plan.AddPlanFeatureData true "PlanFeatures data"
// @Router /api/v1/plans/feature [post]
func CreatePlanFeatures(c *gin.Context) {
	var PlanFeatures models_plan.PlanFeature

	if err := c.ShouldBindJSON(&PlanFeatures); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Create(&PlanFeatures).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "PlanFeatures created successfully",
		"PlanFeatures": PlanFeatures,
	})
}

// GetAllPlanFeatures fetches all PlanFeatures
// @Summary Get all PlanFeatures
// @Description Get all PlanFeatures
// @Tags PlanFeatures
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/v1/plans/feature [get]
func GetAllPlanFeatures(c *gin.Context) {
	var PlanFeatures []models_plan.PlanFeature

	if err := postgres.DB.Preload("PlanFeatureAssociations.Plan").Find(&PlanFeatures).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve PlanFeatures"})
		return
	}

	if len(PlanFeatures) == 0 {
		c.JSON(http.StatusOK, gin.H{"PlanFeatures": []models_plan.PlanFeature{}})
		return
	}

	var featuresResponse []models_plan.PlanFeatureResponse

	for _, feature := range PlanFeatures {
		var featureResponse models_plan.PlanFeatureResponse
		featureResponse.ID = feature.ID
		featureResponse.Name = feature.Name
		featureResponse.Description = feature.Description
		featureResponse.CreatedAt = feature.CreatedAt
		featureResponse.UpdatedAt = feature.UpdatedAt

		for _, association := range feature.PlanFeatureAssociations {

			featureResponse.Plans = append(featureResponse.Plans, models_plan.Plan{
				ID:          association.Plan.ID,
				Name:        association.Plan.Name,
				Description: association.Plan.Description,
				Price:       association.Plan.Price,
				IsActive:    association.Plan.IsActive,
				TrialPeriod: association.Plan.TrialPeriod,
				CreatedAt:   association.Plan.CreatedAt,
				UpdatedAt:   association.Plan.UpdatedAt,
			})
		}

		featuresResponse = append(featuresResponse, featureResponse)
	}

	c.JSON(http.StatusOK, gin.H{"message": "PlanFeatures Found", "PlanFeatures": featuresResponse})
}

// UpdatePlanFeatures updates a PlanFeatures
// @Summary Update a PlanFeatures
// @Description Update a PlanFeatures
// @Tags PlanFeatures
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "PlanFeatures ID"
// @Param input body models_plan.UpdatePlanFeatureData true "PlanFeatures data"
// @Router /api/v1/plans/feature/{id} [put]
func UpdatePlanFeatures(c *gin.Context) {

	id := c.Param("id")
	var PlanFeatures models_plan.PlanFeature

	if err := postgres.DB.First(&PlanFeatures, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PlanFeatures not found"})
		return
	}

	var input models_plan.PlanFeature
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Model(&PlanFeatures).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update PlanFeatures"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "PlanFeatures updated successfully",
		"PlanFeatures": PlanFeatures,
	})
}

// DeletePlanFeatures deletes a PlanFeatures by ID
// @Summary Delete a PlanFeatures
// @Description Delete a PlanFeatures by ID
// @Tags PlanFeatures
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "PlanFeatures ID"
// @Router /api/v1/plans/feature/{id} [delete]
func DeletePlanFeatures(c *gin.Context) {
	id := c.Param("id")
	if err := postgres.DB.Delete(&models_plan.PlanFeature{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete PlanFeatures"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PlanFeatures deleted successfully"})
}
