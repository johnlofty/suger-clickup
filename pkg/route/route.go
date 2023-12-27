package route

import (
	"suger-clickup/pkg/handlers"
	"suger-clickup/pkg/middlewares"
	"suger-clickup/pkg/settings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func BuildRoute(app *fiber.App) {
	// add logger
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	apiV1 := app.Group("/api/v1")
	// need to register before using jwt
	buildUserRoute(apiV1)

	conf := settings.Get()
	jwt := middlewares.NewAuthMiddleware(conf.JWT.Secret)
	apiV1.Use(jwt)

	apiV1.Get("/protected", handlers.Protected)
}

func buildUserRoute(app fiber.Router) {
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
}
