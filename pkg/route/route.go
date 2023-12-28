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

	// ticket
	apiV1.Get("/tickets", h.GetTickets)
	apiV1.Post("/tickets", h.CreateTicket)
	apiV1.Patch("/tickets/:ticket_id/description", h.EditTicketDescription)
	apiV1.Patch("/tickets/:ticket_id/duedate", h.EditTicketDueDate)
	apiV1.Patch("/tickets/:ticket_id/reopen", h.ReopenTicket)
	apiV1.Patch("/tickets/:ticket_id/assignee", h.SetTicketAssignee)
	apiV1.Delete("/tickets/:ticket_id/assignee/:assignee_id", h.DelTicketAssignee)

	// comments
	apiV1.Get("/ticket/:ticket_id/comments", h.GetComments)
	apiV1.Post("/ticket/:ticket_id/comments", h.AddComment)

}
