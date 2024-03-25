package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	db "github.com/patnaikankit/POS-API/config"
	"github.com/patnaikankit/POS-API/models"
)

// create a new cashier
func CreateCashier(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatal("Error while creating new cashier -> ", err)
	}

	if data["name"] == "" || data["passcode"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Both the fields are required!",
			"error":   map[string]interface{}{},
		})
	}

	cashier := models.Cashier{
		Name:      data["name"],
		Passcode:  data["passcode"],
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	db.DB.Create(&cashier)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "New cashier created successfully!",
		"data":    cashier,
	})
}

// edit data of a cashier
func UpdateCashier(c *fiber.Ctx) error {
	cashierID := c.Params("cashierID")
	var cashier models.Cashier

	db.DB.Find(&cashier, "id = ?", cashierID)
	if cashier.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found!",
			"error":   map[string]interface{}{},
		})
	}

	var updateData models.Cashier
	c.BodyParser(&updateData)
	if updateData.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier name is required!",
			"error":   map[string]interface{}{},
		})
	}

	cashier.Name = updateData.Name
	db.DB.Save(&cashier)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Data Updated successfully!",
		"data":    cashier,
	})
}

// retrive all cashiers
func CashierList(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var cashier []models.Cashier
	db.DB.Select("*").Limit(limit).Offset(skip).Find(&cashier).Count(&count)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Data Updated Successfully!",
		"data":    cashier,
	})
}

// get info about a particular cashier
func GetCashierData(c *fiber.Ctx) error {
	cashierID := c.Params("cashierID")

	var cashier models.Cashier
	db.DB.Select("id, name").Where("id=?", cashierID).First(&cashier)
	cashierData := make(map[string]interface{})
	cashierData["cashierID"] = cashier.ID
	cashierData["name"] = cashier.Name

	if cashier.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier NOT Found!",
			"error":   map[string]interface{}{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Data Fetched Succesfully!",
		"data":    cashierData,
	})
}

// delete a cashier
func DeleteCashier(c *fiber.Ctx) error {
	cashierId := c.Params("cashierID")
	var cashier models.Cashier
	db.DB.Where("id=?", cashierId).First(&cashier)

	if cashier.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   map[string]interface{}{},
		})
	}
	db.DB.Where("id = ?", cashierId).Delete(&cashier)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
