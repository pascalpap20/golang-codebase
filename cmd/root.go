package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"example.com/example/config"
	"example.com/example/internal/handler"
	"example.com/example/internal/service"
	"example.com/example/lib/cache"
	"example.com/example/lib/db"
	"example.com/example/lib/transport"
	"git.govtechindonesia.id/inadigital/inalog"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/dubonzi/otelresty"
	"github.com/go-resty/resty/v2"
	//"git.govtechindonesia.id/inadigital/inatrace/integrations/estrace"
	//"git.govtechindonesia.id/inadigital/inatrace/integrations/ddtrace"
)

type Options struct {
	Debug bool   `doc:"Enable debug logging"`
	Host  string `doc:"Hostname to listen on."`
	Port  int    `doc:"Port to listen on." short:"p"`
}

func applyOptions(opts *Options) *config.Config {
	c := config.Get()

	if opts.Host != "" {
		c.Host = opts.Host
	} else if c.Host == "" {
		c.Host = "localhost"
	}

	if opts.Port != 0 {
		c.Port = opts.Port
	} else if c.Port == 0 {
		c.Port = 3000
	}

	return c
}

func Execute() {
	// Then, create the CLI.
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		c := applyOptions(opts)

		inalog.Init(inalog.Cfg{Source: true})

		f := transport.InitFiber(c)

		hooks.OnStart(func() {
			//tp := inatrace.InitTracerDD()
			// OR:
			//tp := estrace.InitTracerES()
			// defer func() {
			// 	if err := tp.Shutdown(context.Background()); err != nil {
			// 		log.Printf("Error shutting down tracer provider: %v", err)
			// 	}
			// }()

			svc := &service.Services{}

			handler.RegisterRoutes(f, svc)

			dbConn, err := db.Open(c)
			if err != nil {
				inalog.Log().Error("Error", slog.Any("error", err))
			}
			svc.DB = dbConn

			cache, err := cache.NewCache(c)
			if err != nil {
				inalog.Log().Error("Error", slog.Any("error", err))
			}
			svc.Cache = cache

			svc.Resty = resty.New()
			opts := []otelresty.Option{otelresty.WithTracerName("example-resty")}
			otelresty.TraceClient(svc.Resty, opts...)

			// Start your server here
			err = f.Listen(fmt.Sprintf("%s:%d", c.Host, c.Port))
			if err != nil {
				inalog.Log().Error("Error", slog.Any("error", err))
			}
		})

		hooks.OnStop(func() {
			// Gracefully shutdown your server here
			f.ShutdownWithTimeout(5 * time.Second)
		})
	})

	rootCmd := cli.Root()
	rootCmd.Use = "example"
	rootCmd.Version = "0.0.1"

	rootCmd.AddCommand(dbInitCmd)
	rootCmd.AddCommand(dbMigrateCmd)
	rootCmd.AddCommand(dbRollbackCmd)
	rootCmd.AddCommand(dbLockCmd)
	rootCmd.AddCommand(dbUnlockCmd)
	rootCmd.AddCommand(dbCreateGoCmd)
	rootCmd.AddCommand(dbCreateSqlCmd)
	rootCmd.AddCommand(dbStatusCmd)
	rootCmd.AddCommand(dbMarkAppliedCmd)

	rootCmd.AddCommand(dbSeedCmd)
	rootCmd.AddCommand(openapiCmd)

	// Run the thing!
	cli.Run()
}
