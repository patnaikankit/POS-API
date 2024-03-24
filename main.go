package main

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/patnaikankit/POS-API/config"
	routes "github.com/patnaikankit/POS-API/routes"
)

func main() {
	db.Connect()

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":3000")
}
