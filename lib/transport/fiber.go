package transport

import (
	"example.com/example/config"
	"github.com/gofiber/fiber/v2"
)

func InitFiber(c *config.Config) *fiber.App {
	f := fiber.New(fiber.Config{})

	return f
}
