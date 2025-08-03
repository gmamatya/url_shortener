package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/gmamatya/url_shortener/database"
	"github.com/gmamatya/url_shortener/helpers"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

// ShortenURL handles the URL shortening request.
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short,omitempty"`
	Expiry      time.Duration `json:"expiry,omitempty"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`       // Remaining requests in the current rate limit window
	XRateLimitReset time.Duration `json:"rate_limit_reset"` // Time until the rate limit resets
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// implement rate limiting logic here
	rdb2 := database.CreateClient(1) // Create a Redis client for rate limiting
	defer rdb2.Close()
	val, err := rdb2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		// Key does not exist, set it with an expiry
		_ = rdb2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := rdb2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Second, // Convert to seconds
			})
		}
	}

	// check if URL is valid
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL format",
		})
	}

	// check for domain validity
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Invalid domain",
		})
	}

	// enforce HTTPS (TLS), SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6] // Generate a new unique ID
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0) // Create a Redis client for the default database
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Custom short URL already exists",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour // Default expiry time
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store URL",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	rdb2.Decr(database.Ctx, c.IP()) // Decrement the rate limit counter

	val, _ = rdb2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := rdb2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute // Convert to seconds

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
