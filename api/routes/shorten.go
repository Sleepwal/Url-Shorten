package routes

import (
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"sleepwalker/url-shorten/helpers"
	"time"
)

type request struct {
	URL            string        `json:"url"`
	CustomShortURL string        `json:"custom_short_url"`
	Expiry         time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShortURL  string        `json:"custom_short_url"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"x_rate_remaining"`   // 速率限制
	XRateLimitReset time.Duration `json:"x_rate_limit_reset"` // 速率限制重置时间
}

func ShortenUrl(c *fiber.Ctx) error {
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	// 速率限制

	// 检查URL是否有效
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL!"})
	}

	// 移除错误
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid domain!"})
	}
	// 强制http，SSL
	body.URL = helpers.ForceHTTP(body.URL)

	return nil
}
