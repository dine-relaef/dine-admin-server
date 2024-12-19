package services

import (
	postgres "menu-server/src/config/database"
	models "menu-server/src/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// @BasePath /api/v1
// CreateOrder handles the creation of a new order
// @Summary Create a new order
// @Description Create a new order with items and options
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrder true "Order data"
// @Router /api/v1/orders [post]
func CreateOrder(c *gin.Context) {
	var input models.CreateOrder

	// Bind the JSON input to the DTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate totals
	var subtotal, tax, serviceFee, total float64
	for _, item := range input.Items {
		var menuItem models.MenuItem
		if err := postgres.DB.First(&menuItem, "id = ?", item.MenuItemID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Menu item not found"})
			return
		}

		var menuItemOptions []models.MenuItemOption
		if err := postgres.DB.Where("menu_item_id = ?", item.MenuItemID).Find(&menuItemOptions).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Menu item options not found"})
			return
		}

		var itemOption models.MenuItemOption
		for _, option := range menuItemOptions {
			itemOptionUUID := *item.ItemOptionID

			if option.ID == itemOptionUUID {
				itemOption = option
				break
			}
		}

		itemTotal := float64(item.Quantity) * itemOption.Price
		subtotal += itemTotal
	}

	tax = subtotal * 0.1         // Example: 10% tax
	serviceFee = subtotal * 0.05 // Example: 5% service fee
	total = subtotal + tax + serviceFee

	// Create the Order
	order := models.Order{
		ID:            uuid.Must(uuid.NewV4()),
		RestaurantID:  input.RestaurantID,
		CustomerEmail: input.CustomerEmail,
		CustomerName:  input.CustomerName,
		CustomerPhone: input.CustomerPhone,
		PaymentType:   input.PaymentType,
		Status:        models.OrderStatusPending,
		OrderType:     models.OrderType(input.OrderType),
		SubTotal:      subtotal,
		Tax:           tax,
		ServiceFee:    serviceFee,
		Total:         total,
		Notes:         input.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := postgres.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create OrderItems
	for _, item := range input.Items {

		var itemOption models.MenuItemOption

		if err := postgres.DB.First(&itemOption, "id = ?", item.ItemOptionID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Menu item option not found"})
			return
		}

		orderItem := models.OrderItem{
			ID:             uuid.Must(uuid.NewV4()),
			OrderID:        order.ID,
			MenuItemID:     item.MenuItemID,
			Quantity:       item.Quantity,
			Price:          itemOption.Price,
			Subtotal:       float64(item.Quantity) * itemOption.Price,
			ItemOptionID:   itemOption.ID,
			ItemOptionName: itemOption.Name,
		}
		if err := postgres.DB.Create(&orderItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order items"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"order":  order,
		"items":  input.Items,
		"status": "Order created successfully",
	})
}

// GetOrder retrieves an order by ID
// @Summary Get order details
// @Description Get detailed information about a specific order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Router /api/v1/orders/{id} [get]
func GetOrder(c *gin.Context) {
	orderID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := postgres.DB.Preload("OrderItems.MenuItem").
		Preload("OrderItems.SelectedOptions").
		First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus updates the status of an order
// @Summary Update order status
// @Description Update the status of a specific order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body models.OrderStatus true "New Status"
// @Router /api/v1/orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	orderID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var statusUpdate struct {
		Status models.OrderStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	var order models.Order
	if err := postgres.DB.First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Update order status
	if err := postgres.DB.Model(&order).Updates(map[string]interface{}{
		"status": statusUpdate.Status,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"order":   order,
	})
}

// ListOrders retrieves orders for a restaurant
// @Summary List restaurant orders
// @Description Get list of orders for a specific restaurant
// @Tags Orders
// @Accept json
// @Produce json
// @Param restaurant_id query string true "Restaurant ID"
// @Param status query string false "Order Status Filter"
// @Router /api/v1/orders [get]
func ListOrders(c *gin.Context) {
	restaurantID, err := uuid.FromString(c.Query("restaurant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
		return
	}

	status := c.Query("status")

	query := postgres.DB.Where("restaurant_id = ?", restaurantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var orders []models.Order
	if err := query.Preload("OrderItems.MenuItem").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// CancelOrder cancels an existing order
// @Summary Cancel order
// @Description Cancel an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Router /api/v1/orders/{id}/cancel [post]
func CancelOrder(c *gin.Context) {
	orderID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := postgres.DB.First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status == models.OrderStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel completed order"})
		return
	}

	if err := postgres.DB.Model(&order).Update("status", models.OrderStatusCancelled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
		"order":   order,
	})
}
