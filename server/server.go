package server

import (
	"github.com/skrevolve/sse/util"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
)

func initMiddlewares(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Cache-Control",
		AllowCredentials: true,
	}))
	app.Use(compress.New(compress.Config{ Level: compress.LevelBestSpeed }))
	app.Use(etag.New())
	app.Use(limiter.New())
}

func Create() *fiber.App {
	app := fiber.New(fiber.Config{
		// Prefork: true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if e, ok := err.(*util.Error); ok {
				return ctx.Status(e.Status).JSON(e)
			} else if e, ok := err.(*fiber.Error); ok {
				return ctx.Status(e.Code).JSON(util.Error{Status: e.Code, ErrCode: 500, Message: e.Message})
			} else {
				return ctx.Status(500).JSON(util.Error{Status: 500, ErrCode: 500, Message: err.Error()})
			}
		},
	})
	initMiddlewares(app)
	return app
}

func Listen(app * fiber.App) error {
	app.Use(func(c *fiber.Ctx) error { return c.SendStatus(404) })
	return app.Listen(":8080")
}