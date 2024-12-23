package services_orders

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	models_dine "menu-server/src/models/dine"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateDineOrder handles the creation of a new DineOrder
// @Summary Create a new DineOrder
// @Description Create a new DineOrder
// @Tags DineOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body models_dine.AddDineOrderData true "DineOrder data"
// @Router /api/v1/orders/dine [post]
func CreateDineOrder(c *gin.Context) {
	var input models_dine.AddDineOrderData

	restaurant_admin_id, _ := c.Get("userID")

	// Bind the JSON input to the DTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var Restaurant models.Restaurant
	if err := postgres.DB.Where("id = ? AND admin_id = ?", input.RestaurantID, restaurant_admin_id).First(&Restaurant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var Plan models.Plan
	if err := postgres.DB.First(&Plan, "id = ?", input.PlanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	var PromoCode models.DinePromoCode
	var DiscountAmount float64
	if input.PromoCode != "" {
		if err := postgres.DB.Where("code = ?", input.PromoCode).First(&PromoCode).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Promo code not found"})
			return
		}

		validFrom, err := time.Parse(time.RFC3339, PromoCode.ValidFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ValidFrom date format"})
			return
		}
		validTo, err := time.Parse(time.RFC3339, PromoCode.ValidTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ValidTo date format"})
			return
		}
		if !PromoCode.IsActive && validFrom.Before(time.Now()) && validTo.After(time.Now()) && PromoCode.MaxUses > 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Promo code is not active"})
			return
		}

		if PromoCode.DiscountType == "percentage" {
			DiscountAmount = Plan.Price * (PromoCode.Discount / 100)
		} else if PromoCode.DiscountType == "amount" {
			DiscountAmount = PromoCode.Discount
		}

		PromoCode.MaxUses -= 1
		if err := postgres.DB.Save(&PromoCode).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update promo code"})
			return
		}
	}

	var DineOrder = models_dine.DineOrder{
		RestaurantID:      input.RestaurantID,
		PlanID:            input.PlanID,
		PromoCode:         input.PromoCode,
		Duration:          input.Duration,
		RestaurantAdminID: Restaurant.AdminID,
		Amount:            Plan.Price,
		DiscountAmount:    DiscountAmount,
		Status:            "pending",
	}

	if err := postgres.DB.Create(&DineOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dine order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Dine order created successfully",
		"dine_order": DineOrder,
	})

}

func GetDineOrders(c *gin.Context) {
	var dineOrders []models_dine.DineOrder

	if err := postgres.DB.Find(&dineOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dine orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dine orders found", "dine_orders": dineOrders})
}

func GetDineOrderByID(c *gin.Context) {
	id := c.Param("id")
	var dineOrder models_dine.DineOrder

	if err := postgres.DB.First(&dineOrder, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dine order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dine_order": dineOrder})
}

func GetDineOrderByUsers(c *gin.Context) {
	var dineOrders []models_dine.DineOrder

	restaurant_admin_id, _ := c.Get("user_id")
	restaurantAdminIDStr, _ := restaurant_admin_id.(string)
	restaurantAdminUUID, _ := uuid.FromString(restaurantAdminIDStr)

	if err := postgres.DB.Where("restaurant_admin_id = ?", restaurantAdminUUID).Find(&dineOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dine orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dine orders found", "dine_orders": dineOrders})
}
