package routes

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestResolveURL(t *testing.T) {
	app := fiber.New()
	app.Get("/:url", func(c *fiber.Ctx) error {
		return c.Redirect("https://www.google.com", fiber.StatusMovedPermanently)
	})

	req := httptest.NewRequest("GET", "/123456", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 301, resp.StatusCode)
	assert.Equal(t, "https://www.google.com", resp.Header.Get("Location"))
}
