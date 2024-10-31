package handler

import (
	"fmt"

	"example.com/example/config"
	"example.com/example/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"

	mLogger "github.com/gofiber/fiber/v2/middleware/logger"
	mRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

func RegisterRoutes(f *fiber.App, svc service.AllServices) huma.API {
	c := config.Get()
	f.Use(mLogger.New())
	f.Use(mRecover.New(mRecover.Config{EnableStackTrace: true}))

	cfg := huma.DefaultConfig("Example API", "0.0.1")
	if c.IsProduction() {
		cfg.DocsPath = ""
	}
	cfg.Servers = []*huma.Server{
		{URL: fmt.Sprintf("http://%s:%d", c.Host, c.Port)},
	}

	api := humafiber.New(f, cfg)

	registerUser(api, svc)

	f.Static("/", "./public")
	f.Use(NotFound)

	return api
}
