package services_payments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	postgres "menu-server/src/config/database"
	"menu-server/src/config/env"
	"menu-server/src/config/payments"
	models_dine "menu-server/src/models/dine"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateDinePayment handles the creation of a new Payment
// @Summary Create a new Payment
// @Description Create a new Payment
// @Tags Payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param order_id path string true "Order ID"
// @Router /api/v1/payments/dine/{order_id} [post]
func CreateDinePayment(c *gin.Context) {
	orderID := c.Param("order_id")

	// Fetch order details from the database
	var Order models_dine.DineOrder
	if err := postgres.DB.First(&Order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Calculate amount after discount
	finalAmount := Order.Amount - Order.DiscountAmount

	// Create a payment link with Razorpay

	client := payments.RazorpayClient
	params := map[string]interface{}{
		"amount":          finalAmount * 100, // Amount in smallest currency unit (paise for INR)
		"currency":        "INR",
		"reference_id":    orderID,
		"description":     "Payment for Order " + orderID,
		"callback_url":    "http://localhost:8080/api/v1/payments/dine/callback",
		"callback_method": "get",
	}
	paymentLink, err := client.PaymentLink.Create(params, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save payment details in the database (optional)
	var Payment = models_dine.DinePayment{
		OrderID:       Order.ID,
		TransactionID: paymentLink["id"].(string),
		Amount:        finalAmount,
		Status:        "pending",
	}
	if err := postgres.DB.Create(&Payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment"})
		return
	}

	// Respond with the payment link
	c.JSON(http.StatusCreated, gin.H{
		"message":      "Payment link created successfully",
		"payment_link": paymentLink["short_url"],
	})
}

// PaymentCallback handles the Razorpay payment callback
// @Summary Razorpay payment callback
// @Description Handle Razorpay payment callback
// @Tags Payments
// @Accept json
// @Produce json
// @Router /api/v1/payments/dine/callback [get]
func PaymentCallback(c *gin.Context) {
	paymentID := c.Query("razorpay_payment_id")
	paymentLinkID := c.Query("razorpay_payment_link_id")
	PaymentReferenceID := c.Query("razorpay_payment_link_reference_id")
	paymentStatus := c.Query("razorpay_payment_link_status")
	signature := c.Query("razorpay_signature")

	// Verify the signature
	expectedSignature := hmac.New(sha256.New, []byte(env.PaymentsVar["RAZORPAY_SECRET_KEY"]))
	expectedSignature.Write([]byte(paymentLinkID + "|" + PaymentReferenceID + "|" + paymentStatus + "|" + paymentID)) // Adjust as per Razorpay documentation
	computedSignature := hex.EncodeToString(expectedSignature.Sum(nil))

	if computedSignature != signature {
		log.Printf("Computed Signature: %s, Received Signature: %s\n", computedSignature, signature)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid payment signature"})
		return
	}

	// Update payment status in the database
	var payment models_dine.DinePayment
	if err := postgres.DB.First(&payment, "transaction_id = ? AND order_id = ?", paymentLinkID, PaymentReferenceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if paymentStatus == "paid" {
		payment.Status = "successful"
		if err := postgres.DB.Save(&payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Payment verified successfully"})
	} else {
		payment.Status = "failed"
		if err := postgres.DB.Save(&payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not successful"})
	}
}
