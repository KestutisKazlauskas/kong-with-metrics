package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.PUT("/events", handlePutRequest)
	return r
}

func TestHandlePutRequest(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Valid Request", func(mt *mtest.T) {
		// Mock the MongoDB client
		client = mt.Client

		// Ensure the database is clean before running the test
		collection := client.Database("eventsdb").Collection("events")
		collection.DeleteMany(context.TODO(), bson.M{})

		router := setupRouter()

		jsonData := `{
			"events": [
				{
					"event": "impression",
					"visitorId": "1234",
					"customerId": "1234",
					"pageUrl": "https://page.url",
					"adId": "1234",
					"timestamp": "2023-10-10T10:10:10.000Z",
					"userAgent": "chrome"
				}
			]
		}`
		req, _ := http.NewRequest("PUT", "/events", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Success"}`, w.Body.String())
	})

	mt.Run("Invalid Timestamp", func(mt *mtest.T) {
		// Mock the MongoDB client
		client = mt.Client

		router := setupRouter()

		jsonData := `{
			"events": [
				{
					"event": "impression",
					"visitorId": "1234",
					"customerId": "1234",
					"pageUrl": "https://page.url",
					"adId": "1234",
					"timestamp": "2023-10-10T10:10:10.000Z",
					"userAgent": "chrome"
				}
			]
		}`
		req, _ := http.NewRequest("PUT", "/events", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Timestamp cannot be in the future")
	})
}
