package routes

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// ResolveURL handles the URL resolution request.
func (h *Handler) ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	ctx := context.Background()

	value, err := h.Rdb0.Get(ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if err := h.Rdb1.Incr(ctx, "counter").Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to increment counter",
		})
	}

	return c.Redirect(value, fiber.StatusMovedPermanently) // Redirect to the original URL
}
