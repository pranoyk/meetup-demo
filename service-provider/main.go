package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"storj.io/common/uuid"
)

type ServiceProvider struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Menu struct {
	ID    string `json:"id" bson:"_id"`
	Gravy string `json:"gravy" bson:"gravy"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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

	// ServiceProvider registration endpoint
	router.POST("/serviceProviders", func(c *gin.Context) {
		// Parse request body
		var serviceProvider ServiceProvider
		if err := c.ShouldBindJSON(&serviceProvider); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate new serviceProvider ID
		ID, err := uuid.New()
		if err != nil {
			return
		}

		serviceProvider.ID = ID.String()

		// Insert serviceProvider into MongoDB
		serviceProvidersCollection := client.Database("serviceProvider-management").Collection("serviceProviders")
		_, err = serviceProvidersCollection.InsertOne(context.Background(), serviceProvider)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, serviceProvider)
	})

	// ServiceProvider authentication endpoint
	router.POST("/serviceProviders/login", func(c *gin.Context) {
		// Parse request body
		var serviceProvider ServiceProvider
		if err := c.ShouldBindJSON(&serviceProvider); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if serviceProvider exists in MongoDB
		serviceProvidersCollection := client.Database("serviceProvider-management").Collection("serviceProviders")
		filter := bson.M{"email": serviceProvider.Email, "password": serviceProvider.Password}
		err := serviceProvidersCollection.FindOne(context.Background(), filter).Decode(&serviceProvider)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, "login successful")
	})

	router.POST("/serviceProviders/menu", func(c *gin.Context) {
		// Parse request body
		var menu Menu
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if serviceProvider exists in MongoDB
		serviceProvidersCollection := client.Database("serviceProvider-management").Collection("menu")
		filter := bson.M{"gravy": menu.Gravy}
		err := serviceProvidersCollection.FindOne(context.Background(), filter).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, menu)
	})

	router.POST("/serviceProviders/validate", func(c *gin.Context) {
		// Parse request body
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(user)
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.Post("http://localhost:8080/users/validate", "application/json", &buf)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if res.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, "valid user")
	})

	// Run Gin server
	if err := router.Run(":8088"); err != nil {
		log.Fatal(err)
	}
}
