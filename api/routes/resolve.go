package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"sleepwalker/url-shorten/database"
)

var StatusHttpRedirect = 301

func ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url") // 参数

	// redis连接
	r := database.CreateClient(0)
	defer func(r *redis.Client) {
		_ = r.Close()
	}(r)

	// 从redis中获取 短链接对应的真实链接
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short Url not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	// 增量redis连接
	rInr := database.CreateClient(1)
	defer func(rInr *redis.Client) {
		_ = rInr.Close()
	}(rInr)

	// 计数器加1
	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, StatusHttpRedirect)
}
