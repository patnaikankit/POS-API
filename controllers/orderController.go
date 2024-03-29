package controllers

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	db "github.com/patnaikankit/POS-API/config"
	"github.com/patnaikankit/POS-API/middlewares"
	"github.com/patnaikankit/POS-API/models"
)

// retrieve all orders
func GetAllOrders(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var order []models.Order

	db.DB.Select("*").Limit(limit).Offset(skip).Find(&order).Count(&count)

	type OrderList struct {
		OrderID        int                `json:"orderID"`
		CashierID      int                `json:"cashiersID"`
		PaymentTypesID int                `json:"paymentTypesID"`
		TotalPrice     int                `json:"totalPrice"`
		TotalPaid      int                `json:"totalPaid"`
		TotalReturn    int                `json:"totalReturn"`
		ReceiptID      string             `json:"receiptID"`
		CreatedAt      time.Time          `json:"createdAt"`
		Payments       models.PaymentType `json:"payment_type"`
		Cashiers       models.Cashier     `json:"cashier"`
	}

	orderResponse := make([]*OrderList, 0)

	for _, v := range order {
		cashier := models.Cashier{}
		db.DB.Where("id = ?", v.CashierID).Find(&cashier)
		paymentType := models.PaymentType{}
		db.DB.Where("id = ?", v.PaymentTypesID).Find(&paymentType)

		orderResponse = append(orderResponse, &OrderList{
			OrderID:        v.ID,
			CashierID:      v.CashierID,
			PaymentTypesID: v.PaymentTypesID,
			TotalPrice:     v.TotalPrice,
			TotalPaid:      v.TotalPaid,
			TotalReturn:    v.TotalReturn,
			ReceiptID:      v.ReceiptID,
			CreatedAt:      v.CreatedAt,
			Payments:       paymentType,
			Cashiers:       cashier,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Orders Fetched Successfully!",
		"data":    orderResponse,
		"meta": map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		},
	})
}

// creating a new order
func NewOrder(c *fiber.Ctx) error {
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

	type Products struct {
		ProductID int `json:"productID"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		PaymentID int        `json:"paymentID"`
		TotalPaid int        `json:"totalPaid"`
		Products  []Products `json:"producs"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Empty Body!",
			"error":   map[string]interface{}{},
		})
	}

	repsonse := make([]*models.ProductResponseOrder, 0)

	var totalInvoicePrice = struct {
		price int
	}{}

	productIDs := ""
	quantities := ""

	for _, v := range body.Products {
		totalPrice := 0
		productIDs = productIDs + "," + strconv.Itoa(v.ProductID)
		quantities = quantities + "," + strconv.Itoa(v.Quantity)

		prod := models.ProductOrder{}
		var discount models.Discount
		db.DB.Table("products").Where("id=?", v.ProductID).First(&prod)
		db.DB.Where("id=?", prod.DiscountID).Find(&discount)
		discountCount := 0

		if discount.Type == "BUY_N" {
			totalPrice = prod.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discountCount = totalPrice - percentage
			totalInvoicePrice.price = totalInvoicePrice.price + discountCount
		}

		if discount.Type == "PERCENT" {
			totalPrice = prod.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discountCount = totalPrice - percentage
			totalInvoicePrice.price = totalInvoicePrice.price + discountCount
		}

		repsonse = append(repsonse, &models.ProductResponseOrder{
			ProductID:        prod.ID,
			Name:             prod.Name,
			Price:            prod.Price,
			Discount:         discount,
			Qty:              v.Quantity,
			TotalNormalPrice: prod.Price,
			TotalFinalPrice:  discountCount,
		})
	}

	// creating a new order in the db
	orderResponse := models.Order{
		CashierID:      1,
		PaymentTypesID: body.PaymentID,
		TotalPrice:     totalInvoicePrice.price,
		TotalPaid:      body.TotalPaid,
		TotalReturn:    body.TotalPaid - totalInvoicePrice.price,
		ReceiptID:      "R000" + strconv.Itoa(rand.Intn(1000)),
		ProductID:      productIDs,
		Quantities:     quantities,
		CreatedAt:      time.Now().Local().UTC(),
		UpdatedAt:      time.Now().Local().UTC(),
	}

	db.DB.Create(&orderResponse)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": true,
		"data": map[string]interface{}{
			"order":    orderResponse,
			"products": repsonse,
		},
	})
}

// retrieve data about a particular order
func OrderData(c *fiber.Ctx) error {
	return nil
}

// calculating the total order
func TotalOrder(c *fiber.Ctx) error {
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

	type Products struct {
		ProductID int `json:"productID"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		Products []Products `json:"producs"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Empty Body!",
			"error":   map[string]interface{}{},
		})
	}

	response := make([]*models.ProductResponseOrder, 0)

	var totalInvoicePrice = struct {
		price int
	}{}

	for _, v := range body.Products {
		totalPrice := 0

		prod := models.ProductOrder{}
		var discount models.Discount
		db.DB.Table("products").Where("id=?", v.ProductID).First(&prod)
		db.DB.Where("id=?", prod.DiscountID).Find(&discount)
		discountCount := 0

		if discount.Type == "BUY_N" {
			totalPrice = prod.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discountCount = totalPrice - percentage
			totalInvoicePrice.price = totalInvoicePrice.price + discountCount
		}

		if discount.Type == "PERCENT" {
			totalPrice = prod.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discountCount = totalPrice - percentage
			totalInvoicePrice.price = totalInvoicePrice.price + discountCount
		}

		response = append(response, &models.ProductResponseOrder{
			ProductID:        prod.ID,
			Name:             prod.Name,
			Price:            prod.Price,
			Discount:         discount,
			Qty:              v.Quantity,
			TotalNormalPrice: prod.Price,
			TotalFinalPrice:  discountCount,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Grand Total Calculated Successfully!",
		"data": map[string]interface{}{
			"subtotal": totalInvoicePrice.price,
			"products": response,
		},
	})
}

func DownloadOrder(c *fiber.Ctx) error {
	return nil
}

func CheckOrder(c *fiber.Ctx) error {
	orderID := c.Params("orderID")
	var order models.Order
	db.DB.Where("id=?", orderID).First(&order)

	if order.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order does not exist",
		})
	}

	if order.IsDownload == 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": map[string]interface{}{
				"isDownload": false,
			},
		})
	}

	if order.IsDownload == 1 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": map[string]interface{}{
				"isDownload": true,
			},
		})
	}

	return nil
}
