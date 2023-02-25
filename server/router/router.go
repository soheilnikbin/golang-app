package router

import (
	"github.com/gin-gonic/gin"
	"sample-app/server/controller"
)

// Router is the struct for configuring the HTTP router
type Router struct {
	//authMiddleware *auth.AuthMiddleware
	authController *controller.AuthController
}

// NewRouter returns a new HTTP router with the configured routes and middleware
func NewRouter(
	//authMiddleware *auth.AuthMiddleware,
	authController *controller.AuthController,
) *gin.Engine {
	// Create a new Gin router
	r := gin.Default()

	// Set up the authentication routes and middleware
	r.POST("/login", authController.Login)
	r.POST("/register", authController.Register)
	//r.Use(authMiddleware.IsAuthenticated())

	return r
}
