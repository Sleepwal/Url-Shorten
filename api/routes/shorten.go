package routes

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"os"
	"sleepwalker/url-shorten/database"
	"sleepwalker/url-shorten/helpers"
	"strconv"
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
	r := database.CreateClient(1)
	defer func(r2 *redis.Client) {
		_ = r2.Close()
	}(r)

	val, err := r.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil { // 该IP不存在
		_ = r.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), time.Minute*30).Err()
	} else { // 该IP已存在
		val, _ = r.Get(database.Ctx, c.IP()).Result()
		count, _ := strconv.Atoi(val)
		if count <= 0 { // 超过api配额，默认是10
			limit, _ := r.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":            "rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// 检查URL是否有效
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL!"})
	}

	// 移除错误
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid domain!"})
	}
	// 强制https，SSL
	body.URL = helpers.ForceHTTPS(body.URL)

	// 速率配额减一，前面都通过才执行
	_ = r.Decr(database.Ctx, c.IP()).Err()

	return nil
}
