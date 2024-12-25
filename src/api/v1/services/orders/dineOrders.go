package services_orders

import (
	postgres "dine-server/src/config/database"
	models_order "dine-server/src/models/orders"
	models_plan "dine-server/src/models/plans"
	models_promoCode "dine-server/src/models/promoCode"
	models_restaurant "dine-server/src/models/restaurants"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateDineOrder handles the creation of a new DineOrder
// @Summary Create a new DineOrder
// @Description Create a new DineOrder
// @Tags Dine Orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body models_order.AddDineOrderData true "DineOrder data"
// @Router /api/v1/orders/dine [post]
func CreateDineOrder(c *gin.Context) error {
	var input models_order.AddDineOrderData

	restaurant_admin_id, _ := c.Get("userID")

	// Bind the JSON input to the DTO
	if err := c.ShouldBindJSON(&input); err != nil {

		return err

	}

	var Restaurant models_restaurant.Restaurant
	if err := postgres.DB.Where("id = ? AND admin_id = ?", input.RestaurantID, restaurant_admin_id).First(&Restaurant).Error; err != nil {
		return err
	}

	var Plan models_plan.Plan
	if err := postgres.DB.First(&Plan, "id = ?", input.PlanID).Error; err != nil {
		return fmt.Errorf("plan not found")
	}

	var PromoCode models_promoCode.DinePromoCode
	var DiscountAmount float64
	if input.PromoCode != "" {
		if err := postgres.DB.Where("code = ?", input.PromoCode).First(&PromoCode).Error; err != nil {
			return fmt.Errorf("promo code not found")
		}

		PlanIDs, err := PromoCode.GetPlanIDs()
		if err != nil {
			return fmt.Errorf("failed to get plan ids")
		}

		planExist := false
		for _, planID := range PlanIDs {
			if planID == Plan.ID {
				planExist = true
				break
			}
		}

		if !planExist {
			return fmt.Errorf("promo code is not applicable to this plan")
		}

		layout := "2006-01-02T15:04:05Z"

		validFrom, err := time.Parse(layout, PromoCode.ValidFrom)
		if err != nil {
			return fmt.Errorf("invalid promo code valid from date")
		}

		validTo, err := time.Parse(layout, PromoCode.ValidTo)
		if err != nil {
			return fmt.Errorf("invalid promo code valid to date")
		}
		if !PromoCode.IsActive || validFrom.Before(time.Now()) || validTo.After(time.Now()) || PromoCode.MaxUses <= 0 {
			return fmt.Errorf("promo code is not applicable")
		}

		if PromoCode.DiscountType == "percentage" {
			DiscountAmount = Plan.Price * (PromoCode.Discount / 100)
		} else if PromoCode.DiscountType == "amount" {
			DiscountAmount = PromoCode.Discount
		}

		if PromoCode.MaxUses > 0 {
			PromoCode.MaxUses -= 1
			if err := postgres.DB.Save(&PromoCode).Error; err != nil {
				return fmt.Errorf("failed to update promo code")
			}
		}
	}

	var DineOrder = models_order.DineOrder{
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
		return fmt.Errorf("failed to create dine order")
	}

	c.Set("orderID", DineOrder.ID.String())

	return nil

}

// GetDineOrders retrieves all dine orders
// @Summary Get all dine orders
// @Description Get all dine orders
// @Tags Dine Orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/v1/orders/dine/all [get]
func GetDineOrders(c *gin.Context) {
	var dineOrders []models_order.DineOrder

	if err := postgres.DB.Find(&dineOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dine orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dine orders found", "dine_orders": dineOrders})
}

// GetDineOrderByID retrieves a dine order by ID
// @Summary Get dine order by ID
// @Description Get dine order by ID
// @Tags Dine Orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Dine Order ID"
// @Router /api/v1/orders/dine/{id} [get]
func GetDineOrderByID(c *gin.Context) {
	id := c.Param("id")
	var dineOrder models_order.DineOrder

	if err := postgres.DB.First(&dineOrder, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dine order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dine_order": dineOrder})
}

// GetDineOrderByUsers retrieves dine orders by restaurant admin
// @Summary Get dine orders by restaurant admin
// @Description Get dine orders by restaurant admin
// @Tags Dine Orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/v1/orders/dine [get]
func GetDineOrderByUsers(c *gin.Context) {
	var dineOrders []models_order.DineOrder

	restaurant_admin_id, _ := c.Get("user_id")
	restaurantAdminIDStr, _ := restaurant_admin_id.(string)
	restaurantAdminUUID, _ := uuid.FromString(restaurantAdminIDStr)

	if err := postgres.DB.Where("restaurant_admin_id = ?", restaurantAdminUUID).Find(&dineOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dine orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dine orders found", "dine_orders": dineOrders})
}
