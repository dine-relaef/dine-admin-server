package controllers

import (
	models "menu-server/src/models"
	postgres "menu-server/src/config/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePayment handles the creation of a new Payment
func CreatePayment(c *gin.Context) {
	var Payment models.Payment

	// Bind JSON to the Payment model
	if err := c.ShouldBindJSON(&Payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save Payment to the postgres
	if err := postgres.DB.Create(&Payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Payment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Payment created successfully",
		"Payment": Payment,
	})
}

// GetPayments retrieves all Payments from the postgres
func GetPayments(c *gin.Context) {
	var Payments []models.Payment

	if err := postgres.DB.Find(&Payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "message" : "Payments Found", "Payments": Payments})
}

// GetPaymentByID retrieves a specific Payment by ID
func GetPaymentByID(c *gin.Context) {
	id := c.Param("id")
	var Payment models.Payment

	if err := postgres.DB.First(&Payment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Payment": Payment})
}

// UpdatePayment updates a Payment's information by ID
func UpdatePayment(c *gin.Context) {
	id := c.Param("id")
	var Payment models.Payment

	if err := postgres.DB.First(&Payment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	var input models.Payment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := postgres.DB.Model(&Payment).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment updated successfully",
		"Payment": Payment,
	})
}

// // DeletePayment removes a Payment by ID
// func DeletePayment(c *gin.Context) {
// 	id := c.Param("id")
// 	if err := postgres.DB.Delete(&models.Payment{}, "id = ?", id).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Payment"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
// }
