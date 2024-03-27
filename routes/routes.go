package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/patnaikankit/POS-API/controllers"
)

func Setup(app *fiber.App) {
	// auth
	app.Post("/cashiers/:cashierID/login", controllers.Login)
	app.Post("/cashiers/:cashierID/logout", controllers.Logout)
	app.Get("/cashiers/:cashierId/passcode", controllers.Passcode)

	// cashier
	app.Get("/cashiers", controllers.CashierList)
	app.Post("/cashiers", controllers.CreateCashier)
	app.Put("/cashiers/:cashierID", controllers.UpdateCashier)
	app.Get("/cashiers/:cashierID", controllers.GetCashierData)
	app.Delete("/cashiers/:cashierID", controllers.DeleteCashier)

	// product
	app.Get("/products", controllers.GetAllProducts)
	app.Post("/products", controllers.AddProduct)
	app.Get("/products/:productID", controllers.GetProductDetails)
	app.Delete("/products/:productID", controllers.DeleteProduct)
	app.Put("/products/:productID", controllers.UpdateProduct)

	// payment
	app.Get("/payments", controllers.GetAllPayments)
	app.Post("/payments", controllers.NewPayment)
	app.Get("/payments/:paymentID", controllers.PaymentData)
	app.Put("/payments/:paymentID", controllers.UpdatePayment)
	app.Delete("/payments/:paymentID", controllers.DeletePayment)
}
