package services_payments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	postgres "menu-server/src/config/database"
	"menu-server/src/config/env"
	"menu-server/src/config/payments"
	models_order "menu-server/src/models/orders"
	models_payment "menu-server/src/models/payments"

	"github.com/gin-gonic/gin"
)

// CreateDinePayment handles the creation of a new Payment
func CreateDinePayment(c *gin.Context) (gin.H, error) {
	orderIDValue, exists := c.Get("orderID")
	if !exists {
		return nil, fmt.Errorf("order id not found")
	}
	orderID := orderIDValue.(string)

	// Fetch order details from the database
	var Order models_order.DineOrder
	if err := postgres.DB.First(&Order, "id = ?", orderID).Error; err != nil {

		return nil, fmt.Errorf("order not found")
	}

	finalAmount := math.Abs(Order.Amount - Order.DiscountAmount)

	// Create a payment link with Razorpay

	if finalAmount == 0 {
		Order.Status = "successful"
		if err := postgres.DB.Save(&Order).Error; err != nil {
			return nil, fmt.Errorf("failed to update payment status")
		}
		return gin.H{
			"message": "Payment verified successfully",
		}, nil
	}
	client := payments.RazorpayClient
	params := map[string]interface{}{
		"amount":          finalAmount * 100, // Amount in smallest currency unit (paise for INR)
		"currency":        "INR",
		"reference_id":    orderID,
		"description":     "Payment for Order " + orderID,
		"callback_url":    "http://localhost:8080/api/v1/workflow/plan/payment-subscription",
		"callback_method": "get",
	}
	paymentLink, err := client.PaymentLink.Create(params, nil)
	if err != nil {

		return nil, fmt.Errorf("failed to create payment link")
	}

	// Save payment details in the database (optional)
	var Payment = models_payment.DinePayment{
		OrderID:       Order.ID,
		TransactionID: paymentLink["id"].(string),
		Amount:        finalAmount,
		Status:        "pending",
	}
	if err := postgres.DB.Create(&Payment).Error; err != nil {

		return nil, fmt.Errorf("failed to create payment")
	}

	// Respond with the payment link
	return gin.H{
		"payment_link": paymentLink["short_url"],
	}, nil
}

// PaymentCallback handles the Razorpay payment callback
// @Summary Razorpay payment callback
// @Description Handle Razorpay payment callback
// @Tags Payments
// @Accept json
// @Produce json
// @Router /api/v1/payments/dine/callback [get]
func PaymentCallback(c *gin.Context) error {
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

		return fmt.Errorf("signature mismatch")
	}

	// Update payment status in the database
	var payment models_payment.DinePayment
	if err := postgres.DB.First(&payment, "transaction_id = ? AND order_id = ?", paymentLinkID, PaymentReferenceID).Error; err != nil {

		return fmt.Errorf("payment not found")
	}

	if paymentStatus == "paid" {
		payment.Status = "successful"
		if err := postgres.DB.Save(&payment).Error; err != nil {

			return fmt.Errorf("failed to update payment status")
		}

		c.Set("paymentID", payment.ID)
		return nil
	} else {
		payment.Status = "failed"
		if err := postgres.DB.Save(&payment).Error; err != nil {

			return fmt.Errorf("failed to update payment status")
		}

		return fmt.Errorf("payment failed")

	}
}
