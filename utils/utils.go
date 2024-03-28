package utils

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	db "github.com/patnaikankit/POS-API/config"
	"github.com/patnaikankit/POS-API/middlewares"
	"github.com/patnaikankit/POS-API/models"
)

func contains(duplicates []int, id int) bool {
	for _, x := range duplicates {
		if x == id {
			return true
		}
	}
	return false
}

func Revenue(c *fiber.Ctx) error {
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

	order := []models.Order{}
	totalRevenue := make([]*models.RevenueResponse, 0)

	Resp1 := models.RevenueResponse{}
	Resp2 := models.RevenueResponse{}

	sum1 := 0
	sum2 := 0

	for _, v := range order {
		if v.PaymentTypesID == 1 {
			payment := models.Payment{}
			paymentTypes := models.PaymentType{}

			db.DB.Where("id=?", 1).First(&paymentTypes)
			db.DB.Where("payment_type_id=?", 1).First(&payment)

			sum1 += v.TotalPaid
			Resp1.Name = payment.Name
			Resp1.Logo = payment.Logo
			Resp1.TotalAmount = sum1
			Resp1.PaymentTypeID = v.PaymentTypesID
		}

		if v.PaymentTypesID == 2 {
			payment := models.Payment{}
			paymentTypes := models.PaymentType{}

			db.DB.Where("id=?", 2).First(&paymentTypes)
			db.DB.Where("payment_type_id=?", 2).First(&payment)

			sum2 += v.TotalPaid
			Resp2.Name = payment.Name
			Resp2.Logo = payment.Logo
			Resp2.TotalAmount = sum2
			Resp2.PaymentTypeID = v.PaymentTypesID
		}
	}

	totalRevenue = append(totalRevenue, &Resp1)
	totalRevenue = append(totalRevenue, &Resp2)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Revenue Calculated Successfully!",
		"data": map[string]interface{}{
			"totalRevenue": sum1 + sum2,
			"paymentTypes": totalRevenue,
		},
	})
}

func Sales(c *fiber.Ctx) error {
	orders := []models.Order{}
	db.DB.Find(&orders)

	sales := make([]*models.SalesResponse, 0)
	totalSales := make([]*models.SalesResponse, 0)

	for _, v := range orders {
		quantities := strings.Split(v.Quantities, ",")
		quantities = quantities[1:]

		products := strings.Split(v.ProductID, ",")
		products = products[1:]

		for i := 0; i < len(products); i++ {
			prod := models.Product{}
			p_id, err1 := strconv.Atoi(products[i])
			quantity, err2 := strconv.Atoi(quantities[i])

			if err1 != nil {
				log.Fatal("Error -> ", err1)
			}

			if err2 != nil {
				log.Fatal("Error -> ", err2)
			}

			db.DB.Where("id", p_id).Find(&prod)
			sales = append(sales, &models.SalesResponse{
				Name:        prod.Name,
				ProductID:   p_id,
				TotalQty:    quantity,
				TotalAmount: quantity * prod.Price,
			})
		}
	}

	duplicates := []int{}
	for _, v := range sales {
		if !contains(duplicates, v.ProductID) {
			duplicates = append(duplicates, v.ProductID)
		}
	}

	quantityArray := []int{}
	for _, v := range duplicates {
		qty := 0
		for _, x := range sales {
			if v == x.ProductID {
				qty = qty + x.TotalQty
			}
		}
		quantityArray = append(quantityArray, qty)
	}

	for i := 0; i < len(duplicates); i++ {
		prod := models.Product{}
		db.DB.Where("id", duplicates[i]).Find(&prod)
		totalSales = append(totalSales, &models.SalesResponse{
			Name:        prod.Name,
			TotalQty:    quantityArray[i],
			TotalAmount: quantityArray[i] * prod.Price,
			ProductID:   duplicates[i],
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Successfully calculated the sales",
		"data": map[string]interface{}{
			"orderProducts": totalSales,
		},
	})
}
