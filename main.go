package main

import (
	"suger-clickup/pkg/models"
	"suger-clickup/pkg/route"
	"suger-clickup/pkg/settings"

	"github.com/gofiber/fiber/v2"
)

func init() {
	settings.Setup()
	models.Setup()

}

func main() {
	app := fiber.New()
	route.BuildRoute(app)
	app.Listen("localhost:3000")
}
