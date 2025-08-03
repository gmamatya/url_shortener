package main

import (
	"context"
	"log"
	"os"

	"github.com/gmamatya/url_shortener/database"
	"github.com/gmamatya/url_shortener/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App, h *routes.Handler) {
	app.Get("/:url", h.ResolveURL)    // resolve redirect
	app.Post("/api/v1", h.ShortenURL) // create short URL
}

func main() {
	ctx := context.Background()
	// --- Create Redis clients once at startup ---
	// db 0: for URL → shortened key storage
	rdb0 := database.CreateClient(ctx, 0)
	// db 1: for lookup / rate-limiting, etc.
	rdb1 := database.CreateClient(ctx, 1)

	// Fiber setup
	app := fiber.New()
	app.Use(logger.New())

	// inject into handlers
	h := routes.NewHandler(rdb0, rdb1)

	// register routes
	setupRoutes(app, h)

	// start
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(port))
}
