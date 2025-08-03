package routes

import (
	"github.com/gmamatya/url_shortener/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// ResolveURL handles the URL resolution request.
func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0) // Create a Redis client for the default database
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	rInr := database.CreateClient(1) // Create a Redis client for the rate limiting database
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, url) // Increment the request count for rate limiting

	return c.Redirect(value, fiber.StatusMovedPermanently) // Redirect to the original URL
}
