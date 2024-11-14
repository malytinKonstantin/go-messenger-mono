package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func IdempotencyMiddleware(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idempotencyKey := c.Get("Idempotency-Key")
		if idempotencyKey == "" {
			return c.Next()
		}

		method := c.Method()
		path := c.Path()

		hash := sha256.Sum256([]byte(idempotencyKey + method + path))
		cacheKey := "idempotency:" + hex.EncodeToString(hash[:])

		val, err := redisClient.Get(context.Background(), cacheKey).Result()
		if err == redis.Nil {
			if err := c.Next(); err != nil {
				return err
			}

			respBody := c.Response().Body()

			expiration := 24 * time.Hour
			if err := redisClient.Set(context.Background(), cacheKey, respBody, expiration).Err(); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Error saving idempotent response")
			}

			return nil
		} else if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Server error")
		}

		c.Response().SetBody([]byte(val))
		c.Response().Header.SetContentType("application/json")
		return nil
	}
}
