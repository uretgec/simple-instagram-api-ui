package main

import (
	"embed"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/peterbourgon/ff/v3"
	"go.uber.org/zap"
)

var ServiceBuild string
var ServiceCommitId string

var ServiceName string = "reweb"
var ServiceVersion string = "1.0.0"

var (
	fs = flag.NewFlagSet(ServiceName, flag.ExitOnError)

	httpAddr    = fs.String("http-addr", "localhost:3001", "http server")
	httpPrefork = fs.Bool("http-prefork", false, "http server prefork option")
	_           = fs.String("env-file", ".env", "env file")
)

//go:embed templates/*
var viewsfs embed.FS

//go:embed assets/*
var assetsfs embed.FS

var sugar *zap.SugaredLogger

func main() {

	// Zap Logger Init
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar = logger.Sugar()
	sugar.Named(ServiceName + ":" + ServiceVersion)

	sugar.Info("service started")

	// Flag Parse with Env
	err := ff.Parse(
		fs, os.Args[1:],
		ff.WithConfigFileFlag("env-file"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix(strings.ToUpper(ServiceName)),
	)

	if err != nil {
		sugar.With("error", err).Fatal("configration error")
	}

	// Fiber Init
	engine := html.NewFileSystem(http.FS(viewsfs), ".tmpl")

	app := fiber.New(fiber.Config{
		Prefork:       *httpPrefork,
		CaseSensitive: true,
		AppName:       ServiceName + " v" + ServiceVersion,
		Views:         engine,
		ViewsLayout:   "templates/layouts/base",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// TODO: may the force be with us (güç bizimle olsun)
			sugar.With("error", err).Fatal("http server stopped")
			return c.Status(fiber.StatusInternalServerError).SendString("502")

		},
		//GETOnly: true,
	})

	// Assets Files: Embeded
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(assetsfs),
		PathPrefix: "assets",
		Browse:     true,
	}))
	//app.Static("/assets", "./assets")

	// Middlewares
	app.Use(recover.New()) // Sometime need to recover all

	// Middlewares
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1: Best
	}))
	app.Use(headerConf())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// Routes
	setupRoutes(app)

	// 404 Middleware
	app.Use(func(c *fiber.Ctx) error {
		// TODO: may the force be with us (güç bizimle olsun)
		return c.Status(fiber.StatusNotFound).SendString("404") // HTTP:404
	})

	go func() {
		err = app.Listen(*httpAddr)
		if err != nil {
			sugar.With("error", err).Fatal("http server stopped")
		}
	}()

	// Listen server quit or something happened and notify channel
	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	<-close

	// Bye bye
	sugar.Info("im shutting down. see you later")

}
