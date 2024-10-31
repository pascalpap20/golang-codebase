package handler

import (
	"fmt"
	"log/slog"

	"example.com/example/config"
	"example.com/example/internal/service"
	"example.com/example/lib/logging"
	"git.govtechindonesia.id/inadigital/inatrace"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"

	mRecover "github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
)

type Handler struct {
	svc service.AllServices
}

func UnwrapFiberUserContextMiddleware(ctx huma.Context, next func(huma.Context)) {
	rctx := inatrace.ContextFrom(ctx.Context())
	ctx = huma.WithContext(ctx, rctx)
	next(ctx)
}

func RegisterRoutes(f *fiber.App, svc service.AllServices) huma.API {
	c := config.Get()
	f.Use(otelfiber.Middleware())
	f.Use(slogfiber.New(logging.Logger()))
	f.Use(mRecover.New(mRecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(_ *fiber.Ctx, e interface{}) {
			stacktraces := logging.GetTrace(1)
			logging.Error("panic", slog.Any("error", e), slog.Any("stacktraces", stacktraces))
		},
	}))

	cfg := huma.DefaultConfig("Example API", "0.0.1")
	if c.IsProduction() {
		cfg.DocsPath = ""
	}
	cfg.Servers = []*huma.Server{
		{URL: fmt.Sprintf("http://%s:%d", c.Host, c.Port)},
	}

	api := humafiber.New(f, cfg)
	api.UseMiddleware(UnwrapFiberUserContextMiddleware)

	h := &Handler{
		svc,
	}

	h.RegisterUser(api)

	f.Static("/", "./public")
	f.Use(NotFound)

	return api
}
