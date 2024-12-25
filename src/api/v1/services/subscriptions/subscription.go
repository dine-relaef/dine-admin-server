package services_subscription

import (
	postgres "dine-server/src/config/database"
	models_order "dine-server/src/models/orders"
	models_payment "dine-server/src/models/payments"
	models_plan "dine-server/src/models/plans"
	models_restaurant "dine-server/src/models/restaurants"
	models_subscription "dine-server/src/models/subscriptions"
	"fmt"
	"log"
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
// @Param input body models_subscription.AddSubscriptionData true "Subscription data"
// @Router /api/v1/subscriptions/{payment_id} [post]
func CreateSubscription(c *gin.Context) (gin.H, error) {
	paymentID, exists := c.Get("paymentID")
	if !exists {

		return nil, fmt.Errorf("payment ID not found")
	}

	var payment models_payment.DinePayment
	if err := postgres.DB.First(&payment, "id = ?", paymentID).Error; err != nil {

		return nil, fmt.Errorf("payment ID not found")
	}

	if payment.Status != "successful" {

		return nil, fmt.Errorf("payment not successful")
	}

	var order models_order.DineOrder
	if err := postgres.DB.First(&order, "id = ?", payment.OrderID).Error; err != nil {

		return nil, fmt.Errorf("order ID not found")
	}

	restaurantAdminID, _ := c.Get("userID")
	restaurantAdminUUID, err := uuid.FromString(restaurantAdminID.(string))

	if err != nil {

		return nil, fmt.Errorf("invalid user ID")
	}

	log.Println(restaurantAdminUUID, " ", order.RestaurantAdminID)
	if restaurantAdminUUID != order.RestaurantAdminID {

		return nil, fmt.Errorf("unauthorized access")
	}

	// Check if subscription already exists for the order
	if err := postgres.DB.Where("order_id = ?", order.ID).First(&models_subscription.Subscription{}).Error; err == nil {

		return nil, fmt.Errorf("subscription already exists for the order")
	}

	var plan models_plan.Plan
	if err := postgres.DB.First(&plan, "id = ?", order.PlanID).Error; err != nil {

		return nil, fmt.Errorf("plan ID not found")
	}
	subscription := models_subscription.Subscription{
		UserID:       order.RestaurantAdminID,
		OrderID:      order.ID,
		RestaurantID: order.RestaurantID,
		PlanID:       order.PlanID,
		Plan:         plan,
		StartDate:    time.Now(),
		EndDate:      CalculateEndDate(time.Now(), order.Duration),
		PaymentID:    payment.ID,
		AutoRenewal:  false,
	}

	if err := postgres.DB.Create(&subscription).Error; err != nil {

		return nil, fmt.Errorf("failed to create subscription")
	}

	if err := postgres.DB.Model(&models_restaurant.Restaurant{}).Where("id = ?", order.RestaurantID).Update("subscription_id", subscription.ID).Error; err != nil {

		return nil, fmt.Errorf("failed to update restaurant subscription ID")
	}
	return gin.H{
		"message":      "Subscription created successfully",
		"subscription": subscription,
	}, nil
}

// GetAllSubscriptions retrieves all subscriptions
// @Summary Get all subscriptions
// @Description Get all subscriptions
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/v1/subscriptions [get]
func GetAllSubscriptions(c *gin.Context) {
	var subscriptions []models_subscription.Subscription

	// Fetch all subscriptions from the database
	if err := postgres.DB.Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

// GetSubscriptionByID retrieves a subscription by its ID
// @Summary Get subscription by ID
// @Description Get subscription by ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Subscription ID"
// @Router /api/v1/subscriptions/{id} [get]
func GetSubscriptionByID(c *gin.Context) {
	id := c.Param("id")
	var subscription models_subscription.Subscription

	// Fetch subscription by ID
	if err := postgres.DB.First(&subscription, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription ID not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription": subscription})
}

// UpdateSubscription updates a subscription by ID
// func UpdateSubscription(c *gin.Context) {
// 	id := c.Param("id")
// 	var subscription models_subscription.Subscription

// 	// Find subscription by ID
// 	if err := postgres.DB.First(&subscription, "id = ?", id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription ID not found"})
// 		return
// 	}

// 	// Bind input data to Subscription struct
// 	var input models_subscription.Subscription
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
// 		return
// 	}

// 	// Update subscription details
// 	if err := postgres.DB.Model(&subscription).Updates(input).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":      "Subscription updated successfully",
// 		"subscription": subscription,
// 	})
// }

// DeleteSubscription deletes a subscription by ID
// func DeleteSubscription(c *gin.Context) {
// 	id := c.Param("id")

// 	// Delete subscription by ID
// 	if err := postgres.DB.Delete(&models_subscription.Subscription{}, "id = ?", id).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
// }

// Helper function to calculate the end date based on duration
func CalculateEndDate(startDate time.Time, duration string) time.Time {
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
