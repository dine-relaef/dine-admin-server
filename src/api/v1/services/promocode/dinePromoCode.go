package services_promocode

import (
	postgres "dine-server/src/config/database"
	models_promoCode "dine-server/src/models/promoCode"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateDinePromoCode handles the creation of a new DinePromoCode
// @Summary Create a new DinePromoCode
// @Description Create a new DinePromoCode
// @Tags Dine Promo Codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body models_promoCode.AddDinePromoCode true "DinePromoCode data"
// @Router /api/v1/promo-code/dine [post]
func CreateDinePromoCode(c *gin.Context) {
	var AddDinePromoCodeData models_promoCode.AddDinePromoCode

	// Bind incoming JSON data to AddDinePromoCodeData struct
	if err := c.ShouldBindJSON(&AddDinePromoCodeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation (optional)
	if len(AddDinePromoCodeData.PlanIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PlanIDs cannot be empty"})
		return
	}

	validFrom := time.Now()
	validTo := validFrom.AddDate(0, 0, AddDinePromoCodeData.Days)
	// Create new promo code entry
	var dinePromoCode = models_promoCode.DinePromoCode{
		Code:         AddDinePromoCodeData.Code,
		Discount:     AddDinePromoCodeData.Discount,
		ValidFrom:    validFrom.Format("2006-01-02T15:04:05Z"),
		ValidTo:      validTo.Format("2006-01-02T15:04:05Z"),
		MaxUses:      AddDinePromoCodeData.MaxUses,
		IsActive:     AddDinePromoCodeData.IsActive,
		DiscountType: AddDinePromoCodeData.DiscountType,
	}
	if err := dinePromoCode.SetPlanIDs(AddDinePromoCodeData.PlanIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Insert the new promo code into the database
	if err := postgres.DB.Create(&dinePromoCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Promo code created successfully",
		"data":    AddDinePromoCodeData,
	})
}
