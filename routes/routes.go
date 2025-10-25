package routes

import (
	"Oratio/handlers"

	"github.com/gin-gonic/gin"
)

//
// THIS IS THE LINE TO FIX
//
func RegisterRoutes(r gin.IRoutes) { // <-- Change *gin.Engine to gin.IRoutes
//
//
//

	// These paths are now relative to the group they are given.
	// So this will become /default/generate
	r.POST("/generate", handlers.GenerateAndStore)
	
	// This will become /default/session
	r.GET("/session", handlers.GetSessionByQuery) // GET /session?id=2
}