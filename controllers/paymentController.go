package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	db "github.com/patnaikankit/POS-API/config"
	"github.com/patnaikankit/POS-API/middlewares"
	"github.com/patnaikankit/POS-API/models"
)

// new payment
func NewPayment(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Error during payment processing -> %v", err)
	}

	if data["name"] == "" || data["type"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Both the fields are Required!",
			"error":   map[string]interface{}{},
		})
	}

	var paymentTypes models.PaymentType
	db.DB.Where("name", data["type"]).First(&paymentTypes)
	payment := models.Payment{
		Name:      data["name"],
		Type:      data["type"],
		PaymentID: int(paymentTypes.ID),
		Logo:      data["logo"],
	}
	db.DB.Create(&payment)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "New Payment Initiated Successfully!",
		"data":    payment,
	})
}

// retrieve all payments
func GetAllPayments(c *fiber.Ctx) error {
	// check if token is present
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Token Not Found!",
		})
	}

	// check if the cashier is authorized
	if err := middlewares.AuthenticateToken(middlewares.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	// pagination
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var payment []models.Payment
	db.DB.Select("id, name, type, payment_type_id, logo, created_at, updated_at").Limit(limit).Offset(skip).Find(&payment).Count(&count)
	meta := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}

	categoriesData := map[string]interface{}{
		"payments": payment,
		"meta":     meta,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data Fetched Successfully!",
		"data":    categoriesData,
	})
}

// retrieve details about a specific payment
func PaymentData(c *fiber.Ctx) error {
	paymentID := c.Params("paymentID")

	// check if token is present
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Token Not Found!",
		})
	}

	// check if the cashier is authorized
	if err := middlewares.AuthenticateToken(middlewares.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	var payment models.Payment
	db.DB.Where("id=?", paymentID).First(&payment)

	if payment.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found!",
			"error":   map[string]interface{}{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Payment Data Fetched Successfully!",
		"data":    payment,
	})
}

// remove a payment
func DeletePayment(c *fiber.Ctx) error {
	paymentID := c.Params("paymentID")

	// check if token is present
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Token Not Found!",
		})
	}

	// check if the cashier is authorized
	if err := middlewares.AuthenticateToken(middlewares.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	var payment models.Payment
	db.DB.First(&payment, paymentID)
	if payment.Name == "" {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "PaymentID doesn't Exist!",
		})
	}

	resposne := db.DB.Delete(&payment)
	if resposne.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Deletion Failed!",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Payment Deleted Successfully!",
	})
}

// update data about a specific payment
func UpdatePayment(c *fiber.Ctx) error {
	paymentID := c.Params("paymentID")
	var totalPayments models.Payment
	db.DB.Find(&totalPayments)

	var payment models.Payment
	db.DB.Find(&payment, "id=?", paymentID)

	var updatedPaymentData models.Payment
	c.BodyParser(&updatedPaymentData)
	if updatedPaymentData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment Name is Required!",
			"error":   map[string]interface{}{},
		})
	}

	var paymentTypeID int
	if updatedPaymentData.Type == "CASH" {
		paymentTypeID = 1
	}
	if updatedPaymentData.Type == "E-WALLET" {
		paymentTypeID = 2
	}

	payment.Name = updatedPaymentData.Name
	payment.Type = updatedPaymentData.Type
	payment.PaymentID = paymentTypeID
	payment.Logo = updatedPaymentData.Logo

	db.DB.Save(&payment)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Paymnet Data Updated Successfully!",
		"data":    payment,
	})
}
