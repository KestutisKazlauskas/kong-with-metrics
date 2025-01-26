package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	Event      string `json:"event" validate:"required,oneof=impression add_tocart page_view"`
	VisitorID  string `json:"visitorId" validate:"required"`
	CustomerID string `json:"customerId" validate:"required"`
	PageURL    string `json:"pageUrl" validate:"required,url"`
	AdID       string `json:"adId" validate:"required"`
	Timestamp  string `json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05.000Z07:00"`
	UserAgent  string `json:"userAgent" validate:"required"`
}

type RequestBody struct {
	Events []Event `json:"events" validate:"required,dive"`
}

var validate *validator.Validate
var client *mongo.Client

func main() {
	validate = validator.New()

	// Get MongoDB connection details from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	mongoUser := os.Getenv("MONGODB_USER")
	mongoPass := os.Getenv("MONGODB_PASS")

	// Connect to MongoDB
	//client := mongo.NewClient(option.ClientOptions{
	//	AppName:  "MyGoApp",
	//	Timeout:  option.DefaultTimeout,
	//})
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//err = client.Connect(ctx, func(c *mongo.Client) {
	//	// Access database and collection
	//	db := c.Database("mydb")
	//	col := db.Collection("mycollection")
	//
	//	// Insert data into MongoDB
	//	_, err = col.InsertOne(ctx, bson.M{
	//		"_id":   data.Name,
	//		"value": data.Value,
	//	})
	//
	//	if err != nil {
	//		log.Printf("Error inserting document: %v", err)
	//		w.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//})
	var err error
	clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(options.Credential{
		Username: mongoUser,
		Password: mongoPass,
	})
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.PUT("/events", handlePutRequest)
	r.Run(":8080")
}

func handlePutRequest(c *gin.Context) {
	var requestBody RequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, event := range requestBody.Events {
		if err := validate.Struct(event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate timestamp
		timestamp, err := time.Parse(time.RFC3339Nano, event.Timestamp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
			return
		}
		if timestamp.After(time.Now().UTC()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Timestamp cannot be in the future"})
			return
		}
	}

	collection := client.Database("eventsdb").Collection("events")
	_, err := collection.InsertMany(context.TODO(), requestBody.Events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
