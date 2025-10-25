package main

import (
	"context"
	"log"

	"Oratio/routes"
	"Oratio/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

// Use the V2 adapter type
var ginLambda *ginadapter.GinLambdaV2

// init is called once when the Lambda cold starts
func init() {
	log.Println("Initializing Lambda function...")
	
	// Initialize database connection
	services.InitDatabase()
	
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	
	// Setup Gin router
	r := gin.Default()

	// --- THIS IS THE FIX ---
	// Create a router group that matches your API Gateway stage name ("default")
	// All requests from API Gateway will start with /default
	stageGroup := r.Group("/default")
	
	// Register all your routes (e.g., /session) onto this group
	// This will make them match /default/session
	routes.RegisterRoutes(stageGroup)
	// -----------------------
	
	// Initialize the V2 adapter
	ginLambda = ginadapter.NewV2(r)
	
	log.Println("Lambda initialization complete")
}

// Handler is the Lambda function handler
// Use the V2 Request and Response types
func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Log the path as received by Lambda
	log.Printf("Received request: %s %s", req.RequestContext.HTTP.Method, req.RequestContext.HTTP.Path)
	
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}