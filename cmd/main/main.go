//Package main specifies that this is an executable command in Go.
//Files under this package are executables to start order microservice and initialize postgres DB.
package main

import (
	"order/microservice"

	"github.com/gin-gonic/gin"
)

//setUpRouter creates a default gin router with appropriate handlers for multiple REST API endpoints for Order service.
func setUpRouter() *gin.Engine {

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	v1 := router.Group("/ordersystem")
	{
		v1.HEAD("/", microservice.Ping)
		v1.POST("/", microservice.POSThandler)
		v1.GET("/get", microservice.GETHandler)
		v1.PUT("/put", microservice.PUTHandler)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	return router
}

//main is called when the executable runs.
//Sets up a web server on port 8080.
func main() {
	router := setUpRouter()
	router.Run()
}
