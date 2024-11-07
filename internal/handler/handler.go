package handler

import (
	"fmt"
	"log/slog"

	"example.com/example/config"
	"example.com/example/internal/service"
	"git.govtechindonesia.id/inadigital/inalog"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"

	mRecover "github.com/gofiber/fiber/v2/middleware/recover"
	mRequestId "github.com/gofiber/fiber/v2/middleware/requestid"
)

type Handler struct {
	svc service.AllServices
}

func UnwrapFiberUserContextMiddleware(ctx huma.Context, next func(huma.Context)) {
	ctx = huma.WithContext(ctx, inalog.WithFiberCtx(ctx.Context()))
	next(ctx)
}

func RegisterRoutes(f *fiber.App, svc service.AllServices) huma.API {
	c := config.Get()
	f.Use(otelfiber.Middleware())
	f.Use(mRequestId.New(mRequestId.Config{ContextKey: inalog.CtxKeyRequestID}))
	f.Use(inalog.NewFiberMiddleware())
	f.Use(mRecover.New(mRecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(_ *fiber.Ctx, e interface{}) {
			inalog.Log().Error("panic", slog.Any("error", e))
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
