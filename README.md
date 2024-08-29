# Go CRUD Server with Loki Logging

This project implements a simple Go-based CRUD server using the Gin framework. The server logs events to Loki for aggregation and visualization. Grafana is used to view and analyze these logs. The project includes:

- **`main.go`**: The Go application that implements CRUD operations and logs messages to Loki.
- **`test.html`**: A basic HTML file for testing CRUD operations via a web interface.
- **Docker**: For running Loki and Grafana containers.
- **Docker Compose**: To manage and configure Loki and Grafana services.

## Dependencies

- **Gin**: Web framework for Go.
- **Promtail Client**: Library for sending logs to Loki.
- **CORS Middleware**: For handling Cross-Origin Resource Sharing.

Install these packages using:

```sh
go get github.com/gin-contrib/cors
go get github.com/gin-gonic/gin
go get github.com/afiskon/promtail-client/promtail
```

## `main.go`

### Overview

`main.go` is the core of the CRUD server application. It initializes a Loki client, sets up an HTTP server using Gin, and logs all HTTP requests and responses to Loki.

### Code Explanation

```go
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

	// Define CRUD routes
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
```

### Explanation

- **`initLoki`**: Configures and initializes the Loki client with the necessary settings.
- **`logToLoki`**: Sends formatted log messages to Loki.
- **`main`**: Sets up the Gin router, logs HTTP requests and responses, and defines CRUD routes.
- **CRUD Operations**: Handlers for creating, reading, updating, and deleting items, with logging to Loki.

## `test.html`

### Overview

`test.html` provides a user interface for testing CRUD operations. It includes buttons to trigger HTTP requests to the Go server and displays the server's JSON responses.

### Code

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CRUD Test Page</title>
    <script>
        async function testCRUD(method, id = '') {
            const url = `http://<ec2-ip>:5000/item${id ? `/${id}` : ''}`;
            const response = await fetch(url, { method });
            const result = await response.json();
            document.getElementById('result').textContent = JSON.stringify(result, null, 2);
        }
    </script>
</head>
<body>
    <h1>CRUD Test Page</h1>
    <button onclick="testCRUD('GET', '123')">Get Item</button>
    <button onclick="testCRUD('POST')">Create Item</button>
    <button onclick="testCRUD('PUT', '123')">Update Item</button>
    <button onclick="testCRUD('DELETE', '123')">Delete Item</button>
    <pre id="result"></pre>
</body>
</html>
```

### Explanation

- **JavaScript Functions**: `testCRUD` performs fetch requests to the server based on the provided HTTP method and ID.
- **Buttons**: Trigger CRUD operations.
- **Result Display**: Shows the JSON response from the server.



## Docker Compose File

### `docker-compose.yml`

This file defines the configuration for running Loki and Grafana using Docker Compose.

```yaml
version: '3'
services:
  loki:
    image: grafana/loki:latest
    container_name: loki
    ports:
      - "3100:3100"
    networks:
      - my-network

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    networks:
      - my-network

networks:
  my-network:
    driver: bridge
```

### Explanation

- **Loki**: Log aggregation system, exposed on port 3100.
- **Grafana**: Visualization tool, exposed on port 3000.
- **Network**: Custom Docker network `my-network` for communication between containers.

### Running Docker Compose

Start the services with:

```sh
docker-compose up -d
```

This will launch Loki and Grafana in detached mode.

### Run Server
To run the go server use 

```sh
go run main.go
```

### Verify the routes

```sh
http://<ec2-ip>:5000/test
```

## Grafana Setup with Loki

### 1. Access Grafana

Open your browser and go to `http://<ec2-ip>:3000`.

### 2. Add Loki Data Source

1. **Log In**: Default credentials are `admin` / `admin`.
2. **Add Data Source**:
   - Go to the gear icon (⚙️) > "Data Sources".
   - Click "Add data source".
   - Select "Loki" from the list.
   - Set the URL to `http://loki:3100` (for Docker Compose) or `http://localhost:3100` if not using Docker Compose.
   - Click "Save & Test" to verify the connection.

### 3. Create a Dashboard

1. **Create Dashboard**:
   - Click the "+" icon > "Dashboard" > "Add new panel".
   - Set the data source to "Loki".
   - Choose Table in visualization in dashboard
   - From Lab-filters drop down select the `label` we created 

2. **Customize and Save**:
   - Adjust panel settings as needed.
   - Click "Apply" to add the panel to your dashboard.
   - Save the dashboard by clicking the disk icon (💾).

   ![alt text](./go-crud-server/images/image.png)
