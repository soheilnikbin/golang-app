package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"log"
	"sample-app/auth"
	"sample-app/db"
	"sample-app/server/controller"
	"sample-app/server/router"
)

func main() {
	// Set up the database connection
	dbConn, err := db.ConnectToDatabase("host=localhost port=5432 user=postgres dbname=postgres password=member sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	defer dbConn.Close()

	// Create a Firebase app instance
	opt := option.WithCredentialsFile("./service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}

	// Create a Firebase auth client instance
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Firebase auth client: %v", err)
	}

	// Run database migrations
	dbConn.AutoMigrate(&auth.User{})

	// Set up the authentication service and middleware
	authService := &auth.AuthService{
		DB:       dbConn,
		FireAuth: authClient,
	}
	authController := controller.NewAuthController(authService)

	// Set up the HTTP router
	r := router.NewRouter(authController)

	// Start the server
	r.Run(":8080")
}
