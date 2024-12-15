package controllers

import (
	models "menu-server/src/models"
	postgres "menu-server/src/config/database"
	"github.com/gin-gonic/gin"
	"net/http"
	// "strconv"
)

// CreateSubscription handles creating a new subscription
func CreateSubscription(c *gin.Context) {
	var subscription models.Subscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save subscription to postgres
	if err := postgres.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}


	c.JSON(http.StatusCreated, gin.H{
		"message": "subscriptions Paln created successfully",
		"subscription": subscription,
	})
}

// GetAllSubscriptions retrieves all subscriptions
func GetAllSubscriptions(c *gin.Context) {
	var subscription []models.Subscription

	if err := postgres.DB.Find(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": subscription})
}

// GetSubscriptionByID retrieves a subscription by its ID
func GetSubscriptionByID(c *gin.Context) {
	id := c.Param("id")
	var subscription models.Subscription

	if err := postgres.DB.First(&subscription, "id = ?", id).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription Plan ID not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription": subscription})
}

// UpdateSubscription updates a subscription by ID
func UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var subscription models.Subscription

	if err := postgres.DB.First(&subscription, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription plan ID not found"})
		return
	}

	var input models.Subscription
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Model(&subscription).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription Plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subscription Plan updated successfully",
		"subscription": subscription,
	})
}

// DeleteSubscription deletes a subscription by ID
func DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	var subscription models.Subscription
	if err := postgres.DB.Delete(&subscription, "id = ?", id).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid subscription ID"})
	}


	c.JSON(http.StatusOK, gin.H{"message" : "subscription plan deleted successfully"})
}
