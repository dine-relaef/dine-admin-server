package controllers

import (
	models "menu-server/src/models"
	postgres "menu-server/src/config/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateOrder handles the creation of a new Order
func CreateOrder(c *gin.Context) {
	var Order models.Order

	// Bind JSON to the Order model
	if err := c.ShouldBindJSON(&Order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save Order to the postgres
	if err := postgres.DB.Create(&Order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"Order": Order,
	})
}

// GetOrders retrieves all Orders from the postgres
func GetOrders(c *gin.Context) {
	var Orders []models.Order

	if err := postgres.DB.Find(&Orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "message" : "Orders Found", "Orders": Orders})
}

// GetOrderByID retrieves a specific Order by ID
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var Order models.Order

	if err := postgres.DB.First(&Order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Order": Order})
}

// UpdateOrder updates a Order's information by ID
func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var Order models.Order

	if err := postgres.DB.First(&Order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var input models.Order
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Model(&Order).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order updated successfully",
		"Order": Order,
	})
}

// DeleteOrder removes a Order by ID
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := postgres.DB.Delete(&models.Order{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
