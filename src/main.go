package main

import (
	"log"
	postgres "menu-server/src/config/database"
	users "menu-server/src/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	postgres.InitDB()
	r := gin.Default()

	// Setup routes
	// setupRoutes(r, db)

	//health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Jai Mata Di",
			"status":  "ok",
		})
	})

	setupRoutes(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}

func createUser(c *gin.Context) {
	var user users.User

	// Bind JSON to the User model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save user to the database
	if err := postgres.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func getUsers(c *gin.Context) {
	var users []users.User

	// Retrieve all users from the database
	if err := postgres.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func setupRoutes(r *gin.Engine) {
	userRoutes(r,"/users")
}

func userRoutes(r *gin.Engine, prefix string) {
	r.POST(prefix, createUser)
	r.GET(prefix, getUsers)
}
