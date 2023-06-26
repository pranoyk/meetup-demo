package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"storj.io/common/uuid"
)

type User struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func main() {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Set up Gin router
	router := gin.Default()

	// User registration endpoint
	router.POST("/users", func(c *gin.Context) {
		// Parse request body
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate new user ID
		ID, err := uuid.New()
		if err != nil {
			return
		}

		user.ID = ID.String()

		// Insert user into MongoDB
		usersCollection := client.Database("user-management").Collection("users")
		_, err = usersCollection.InsertOne(context.Background(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	})

	// User authentication endpoint
	router.POST("/users/login", func(c *gin.Context) {
		// Parse request body
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user exists in MongoDB
		usersCollection := client.Database("user-management").Collection("users")
		filter := bson.M{"email": user.Email, "password": user.Password}
		err := usersCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	router.POST("/users/validate", func(c *gin.Context) {
		// Parse request body
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user exists in MongoDB
		usersCollection := client.Database("user-management").Collection("users")
		filter := bson.M{"email": user.Email, "name": user.Name}
		err := usersCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or name"})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	// Run Gin server
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
