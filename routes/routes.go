package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/patnaikankit/POS-API/controllers"
)

func Setup(app *fiber.App) {
	// auth
	app.Post("cashiers/:cashierId/login", controllers.CreateCashier)

	// cashier
	app.Get("/cashiers", controllers.CashierList)
	app.Post("/cashiers", controllers.CreateCashier)
	app.Put("/cashiers/:cashierID", controllers.UpdateCashier)
	app.Get("/cashiers/:cashierID", controllers.GetCashierData)
	app.Delete("/cashiers/:cashierID", controllers.DeleteCashier)

	// product

	// payment
}
