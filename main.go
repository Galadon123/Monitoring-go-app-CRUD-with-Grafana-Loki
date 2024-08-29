package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var lokiClient promtail.Client

func initLoki() {
	var err error
	config := promtail.ClientConfig{
		PushURL:            "http://localhost:3100/loki/api/v1/push",
		Labels:             "{app=\"go-crud-server\"}",
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.ERROR,
	}

	lokiClient, err = promtail.NewClientProto(config)
	if err != nil {
		log.Fatalf("Failed to create Loki client: %v", err)
	}
}

func logToLoki(message string) {
	lokiClient.Infof(message)
}

func main() {
	initLoki()
	defer lokiClient.Shutdown()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		// Log the request to Loki
		logToLoki(fmt.Sprintf("Received %s request for %s", c.Request.Method, c.Request.URL.Path))

		c.Next()

		// Log the response to Loki
		logToLoki(fmt.Sprintf("Sent response with status %d for %s request to %s", c.Writer.Status(), c.Request.Method, c.Request.URL.Path))
	})

	// Define your CRUD routes
	r.POST("/item", createItem)
	r.GET("/item/:id", getItem)
	r.PUT("/item/:id", updateItem)
	r.DELETE("/item/:id", deleteItem)

	r.Use(cors.Default())
	r.StaticFile("/test", "./test.html")

	r.Run(":5000")
}

// CRUD operations
func createItem(c *gin.Context) {
	logToLoki("Creating new item")
	c.JSON(http.StatusCreated, gin.H{"message": "Item created successfully"})
}

func getItem(c *gin.Context) {
	id := c.Param("id")
	logToLoki(fmt.Sprintf("Getting item with ID: %s", id))
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Get item with ID: %s", id)})
}

func updateItem(c *gin.Context) {
	id := c.Param("id")
	logToLoki(fmt.Sprintf("Updating item with ID: %s", id))
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Item with ID %s updated successfully", id)})
}

func deleteItem(c *gin.Context) {
	id := c.Param("id")
	logToLoki(fmt.Sprintf("Deleting item with ID: %s", id))
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Item with ID %s deleted successfully", id)})
}