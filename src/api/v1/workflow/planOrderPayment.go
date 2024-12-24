package workflow

import (
	"log"
	services_orders "menu-server/src/api/v1/services/orders"
	services_payments "menu-server/src/api/v1/services/payments"
	services_subscription "menu-server/src/api/v1/services/subscriptions"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PlanOrderPayment is a workflow function that creates a dine order and payment
// @Summary Create a dine order and payment
// @Description Create a dine order and payment
// @Tags Workflow
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body models_order.AddDineOrderData true "DineOrder data"
// @Router /api/v1/workflow/plan/order-payment [post]
func PlanOrderPayment(c *gin.Context) {

	// Step 1: Create the dine order
	// Call the service function to create the dine order
	err := services_orders.CreateDineOrder(c)
	if err != nil { // Check if order creation failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return // If error occurred, stop execution
	}

	// Step 2: After successful dine order creation, create the dine payment
	res, err := services_payments.CreateDinePayment(c)
	if err != nil { // Check if payment creation failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return // If error occurred, stop execution
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order and payment created successfully",
		"payment": res,
	})
}

// VerifyPaymentAndSubscription is a workflow function that verifies a payment and creates a dine subscription
// @Summary Verify a payment and create a dine subscription
// @Description Verify a payment and create a dine subscription
// @Tags Workflow
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/v1/workflow/plan/payment-subscription [post]
func VerifyPaymentAndSubscription(c *gin.Context) {
	// Step 1: Verify the payment
	// Call the service function to verify the payment
	paymentID := c.Query("razorpay_payment_id")
	paymentLinkID := c.Query("razorpay_payment_link_id")
	PaymentReferenceID := c.Query("razorpay_payment_link_reference_id")
	paymentStatus := c.Query("razorpay_payment_link_status")
	signature := c.Query("razorpay_signature")

	log.Println("Payment ID: ", paymentID)
	log.Println("Payment Link ID: ", paymentLinkID)
	log.Println("Payment Reference ID: ", PaymentReferenceID)
	log.Println("Payment Status: ", paymentStatus)
	log.Println("Signature: ", signature)

	err := services_payments.PaymentCallback(c)
	if err != nil { // Check if payment verification failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return // If error occurred, stop execution
	}
	// Step 2: After successful payment verification, create the dine subscription
	res, err := services_subscription.CreateSubscription(c)
	if err != nil { // Check if subscription creation failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return // If error occurred, stop execution
	}

	c.JSON(http.StatusCreated, res)
}
