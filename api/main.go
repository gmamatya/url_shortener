package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gmamatya/url_shortener/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	// Initialize routes here
	// This function will set up the HTTP routes for the URL shortening service.
	// It will include handlers for shortening URLs, retrieving shortened URLs,
	// and any other necessary endpoints.
	app.Get("/:url", routes.ResolveURL)    // Handler for resolving shortened URLs
	app.Post("/api/v1", routes.ShortenURL) // Handler for shortening URLs
}

func main() {
	err := godotenv.Load() // Load environment variables from .env file

	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()

	app.Use(logger.New()) // Middleware for logging requests

	// Set up routes
	setupRoutes(app)

	// Start the server
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
