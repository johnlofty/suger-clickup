package main

import (
	"suger-clickup/pkg/clients"
	"suger-clickup/pkg/dao"
	"suger-clickup/pkg/models"
	"suger-clickup/pkg/route"
	"suger-clickup/pkg/services"
	"suger-clickup/pkg/settings"

	"github.com/gofiber/fiber/v2"
)

func init() {
	settings.Setup()
	clients.Setup()
	models.Setup()

}

func main() {
	app := fiber.New()
	conf := settings.Get()

	s := services.NewService(
		dao.NewDBDao(models.GetDB()),
		clients.NewClickupHandler(
			conf.ClickupConfig.Secret,
			conf.ClickupConfig.ListId),
	)
	route.BuildRoute(app, s)
	app.Listen("localhost:3000")
}
