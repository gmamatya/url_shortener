package routes

import (
	"os"
	"strconv"
	"time"

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

func (h *Handler) ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	ctx := c.Context()

	// implement rate limiting logic here
	rdb2 := h.Rdb1 // Use the second Redis client for rate limiting
	val, err := rdb2.Get(ctx, c.IP()).Result()
	if err == redis.Nil {
		// Key does not exist, set it with an expiry
		quota, err := strconv.ParseInt(os.Getenv("API_QUOTA"), 10, 64)
		if err != nil {
			quota = 10
		}
		if err := rdb2.Set(ctx, c.IP(), quota, 30*60*time.Second).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set rate limit"})
		}
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get rate limit"})
	} else {
		valInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse rate limit"})
		}
		if valInt <= 0 {
			limit, err := rdb2.TTL(ctx, c.IP()).Result()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get rate limit expiration"})
			}
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit.Seconds(),
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
	if !helpers.IsServiceDomain(body.URL) {
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

	val, err = h.Rdb0.Get(ctx, id).Result()
	if err != redis.Nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Custom short URL already exists",
		})
	}

	if body.Expiry == 0 {
		expiry, err := strconv.ParseInt(os.Getenv("DEFAULT_EXPIRY"), 10, 64)
		if err != nil {
			expiry = 24
		}
		body.Expiry = time.Duration(expiry) * time.Hour
	}

	if err := h.Rdb0.Set(ctx, id, body.URL, body.Expiry).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store URL",
		})
	}

	resp := response{
		URL:         body.URL,
		CustomShort: "",
		Expiry:      body.Expiry,
	}

	if err := rdb2.Decr(ctx, c.IP()).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update rate limit"})
	}

	val, err = rdb2.Get(ctx, c.IP()).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get rate limit"})
	}
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, err := rdb2.TTL(ctx, c.IP()).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get rate limit expiration"})
	}
	resp.XRateLimitReset = time.Duration(ttl.Minutes()) * time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
