package services

import (
	"log"
	postgres "menu-server/src/config/database"
	"menu-server/src/models"
	models_dine "menu-server/src/models/dine"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateSubscription handles creating a new subscription
// @Summary Create a new subscription
// @Description Create a new subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payment_id path string true "Payment ID"
// @Param input body models.AddSubscriptionData true "Subscription data"
// @Router /api/v1/subscriptions/{payment_id} [post]
func CreateSubscription(c *gin.Context) {
	paymentID := c.Param("payment_id")

	if _, err := uuid.FromString(paymentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID format"})
		return
	}

	var payment models_dine.DinePayment
	if err := postgres.DB.First(&payment, "id = ?", paymentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment ID not found"})
		return
	}

	if payment.Status != "successful" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Payment not successful"})
		return
	}

	var order models_dine.DineOrder
	if err := postgres.DB.First(&order, "id = ?", payment.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order ID not found"})
		return
	}

	restaurantAdminID, _ := c.Get("userID")
	restaurantAdminUUID, err := uuid.FromString(restaurantAdminID.(string))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Println(restaurantAdminUUID, " ", order.RestaurantAdminID)
	if restaurantAdminUUID != order.RestaurantAdminID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input models.AddSubscriptionData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// Check if subscription already exists for the order
	if err := postgres.DB.Where("order_id = ?", order.ID).First(&models.Subscription{}).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Subscription for this order already exists"})
		return
	}

	subscription := models.Subscription{
		UserID:       order.RestaurantAdminID,
		OrderID:      order.ID,
		RestaurantID: order.RestaurantID,
		PlanID:       order.PlanID,
		StartDate:    time.Now(),
		EndDate:      calculateEndDate(time.Now(), order.Duration),
		PaymentID:    payment.ID,
		AutoRenewal:  input.AutoRenewal,
	}

	if err := postgres.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Model(&models.Restaurant{}).Where("id = ?", order.RestaurantID).Update("subscription_id", subscription.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant subscription"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":      "Subscription created successfully",
		"subscription": subscription,
	})
}

// GetAllSubscriptions retrieves all subscriptions
func GetAllSubscriptions(c *gin.Context) {
	var subscriptions []models.Subscription

	// Fetch all subscriptions from the database
	if err := postgres.DB.Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

// GetSubscriptionByID retrieves a subscription by its ID
func GetSubscriptionByID(c *gin.Context) {
	id := c.Param("id")
	var subscription models.Subscription

	// Fetch subscription by ID
	if err := postgres.DB.First(&subscription, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription ID not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription": subscription})
}

// UpdateSubscription updates a subscription by ID
func UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var subscription models.Subscription

	// Find subscription by ID
	if err := postgres.DB.First(&subscription, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription ID not found"})
		return
	}

	// Bind input data to Subscription struct
	var input models.Subscription
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// Update subscription details
	if err := postgres.DB.Model(&subscription).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Subscription updated successfully",
		"subscription": subscription,
	})
}

// DeleteSubscription deletes a subscription by ID
func DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	// Delete subscription by ID
	if err := postgres.DB.Delete(&models.Subscription{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

// Helper function to calculate the end date based on duration
func calculateEndDate(startDate time.Time, duration string) time.Time {
	// Example logic to parse and calculate the end date
	// Implement proper duration parsing logic
	switch duration {
	case "1M":
		return startDate.AddDate(0, 1, 0)
	case "6M":
		return startDate.AddDate(0, 6, 0)
	case "1Y":
		return startDate.AddDate(1, 0, 0)
	default:
		// Default to 1 month if the duration is invalid
		return startDate.AddDate(0, 1, 0)
	}
}

// Placeholder function to calculate subscription price
func calculatePrice(planID uuid.UUID, promoCode string) float64 {
	// Implement your logic to determine the price based on the plan and promo code
	return 100.0 // Example: return a default price
}
