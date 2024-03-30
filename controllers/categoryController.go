package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	db "github.com/patnaikankit/POS-API/config"
	"github.com/patnaikankit/POS-API/middlewares"
	"github.com/patnaikankit/POS-API/models"
)

type Category struct {
	ID   int    `json:"categoryID"`
	Name string `json:"name"`
}

func AddCategory(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Error while parsing -> %v", err)
	}

	if data["name"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category Name is Required!",
			"error":   map[string]interface{}{},
		})
	}

	category := models.Category{
		Name: data["name"],
	}

	db.DB.Create(&category)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "New Category Added!",
		"error":   category,
	})
}

func GetAllCatgory(c *fiber.Ctx) error {
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
	var category []Category
	db.DB.Select("id ,name").Limit(limit).Offset(skip).Find(&category).Count(&count)
	meta := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}

	categoryData := map[string]interface{}{
		"categories": category,
		"meta":       meta,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "All Categories Fetched!",
		"data":    categoryData,
	})
}

func CategoryData(c *fiber.Ctx) error {
	categoryID := c.Params("categoryID")

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

	var category models.Category
	db.DB.Select("id, name").Where("id=?", categoryID).First(&category)
	categoryData := make(map[string]interface{})
	categoryData["categoryID"] = category.ID
	categoryData["name"] = category.Name

	if category.Name == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "No Category Found!",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Data Fetched Successfully!",
		"data":    categoryData,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
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

	categoryID := c.Params("categoryID")
	var category models.Category

	db.DB.Find(&category, "id=?", categoryID)

	if category.Name == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "No Category Found!",
		})
	}

	var UpdateCashierData models.Category
	c.BodyParser(&UpdateCashierData)
	if UpdateCashierData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category Name is Required!",
			"error":   map[string]interface{}{},
		})
	}

	category.Name = UpdateCashierData.Name
	db.DB.Save(&category)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Category Updated Successfully!",
		"data":    category,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
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

	categoryID := c.Params("categoryID")
	var category models.Category
	db.DB.Where("id=?", categoryID).First(&category)

	if category.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found!",
			"error":   map[string]interface{}{},
		})
	}

	db.DB.Where("id=?", category.ID).Delete(&category)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Category Deleted Successfully!",
	})
}
