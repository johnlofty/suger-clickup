package route

import (
	"suger-clickup/pkg/handlers"
	"suger-clickup/pkg/middlewares"
	"suger-clickup/pkg/services"
	"suger-clickup/pkg/settings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func BuildRoute(app *fiber.App, s *services.Service) {
	// add logger
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	h := handlers.NewHandler(s)

	apiV1 := app.Group("/api/v1")
	apiV1.Post("/register", h.Register)
	apiV1.Post("/login", h.Login)

	conf := settings.Get()
	jwt := middlewares.NewAuthMiddleware(conf.JWT.Secret)
	apiV1.Use(jwt)

	// user
	apiV1.Patch("/user", h.UpdateUser)
	// org
	apiV1.Post("/org", h.CreateOrg)
	// notification
	apiV1.Get("/org/:org_id/notification", h.GetOrgNotification)
	apiV1.Patch("/org/:org_id/notification", h.ModOrgNotification)

	// ticket
	apiV1.Get("/tickets", h.ListTickets)
	apiV1.Get("/tickets/:ticket_id", h.GetTicket)
	apiV1.Post("/tickets", h.CreateTicket)
	apiV1.Patch("/tickets/:ticket_id", h.EditTicket)
	// comments
	apiV1.Post("/tickets/:ticket_id/comments", h.AddComment)
	// TODO update status, validate status
	apiV1.Patch("/tickets/:ticket_id/status", h.ReopenTicket)
}
